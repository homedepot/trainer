package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"bytes"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/homedepot/trainer/actions"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/security"
	"github.com/homedepot/trainer/structs/expected"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/homedepot/trainer/structs/planaction"
	"github.com/juju/loggo"
	"github.com/mitchellh/mapstructure"
	"github.com/mohae/deepcopy"
	"os"
	"reflect"
	"strconv"
	"text/template"
	"time"
)

type test struct {
	tst        *plan.Plan
	processing bool
	action     actions.Action
	bgActions  []actions.Action // actions that could be, but not necessarily are, backgrounded.
}

var tst test

type SatisfyGroup struct {
	name   string
	action []*planaction.PlanAction
}

func LaunchTest(cfg *config.Config, p string) error {
	if tst.tst != nil {
		return errors.New("test already in progress")
	}
	pl, err := cfg.FindPlan(p)
	if err != nil {
		return err
	}
	pln := deepcopy.Copy(pl)
	tst.tst = pln.(*plan.Plan)

	err = tst.tst.Reset()
	if err != nil {
		return err
	}
	return nil
}

func RemoveTest(u uuid.UUID) error {
	if tst.tst == nil {
		return errors.New("no test to remove")
	}
	return nil
}

func ProcessTests(kill chan *bool) {
	logger := loggo.GetLogger("default")
	tst.processing = true
	var ex, nex expected.Expected
	defer func() {
		tst.processing = false
	}()
	if tst.tst == nil {
		return
	}
	if tst.tst.StopVar != "" {
		sv, ok := tst.tst.State.Variables[tst.tst.StopVar]
		if ok {
			svb, ok1 := sv.(bool)
			if ok1 {
				if svb == true {
					tst.tst.State.States[len(tst.tst.State.States)-1].Status = "stopped"
					return
				}
			}
		}
	}
	if tst.tst.State.Err != nil {
		return
	}
	txn, err := tst.tst.FindTransaction(tst.tst.State.Transaction)
	if err != nil {
		logger.Warningf("invalid transaction in state, not processing")
		tst.tst.State.Err = errors.New("invalid transaction in state")
		return
	}
	ex = txn.OnExpected
	nex = txn.OnUnexpected

	// would like to make this a part of Transaction, but can't because of import loops.
	// probably means I still designed something wrong, but meh.
	// sometimes go's composition just gets in the way
	var urlres *actions.ExecuteResult
	logger.Tracef("TxnActionIdx: %v len(txn.InitAction) %v", tst.tst.State.TxnActionIdx, len(txn.InitAction))
	var ctx *actions.QueueContext
	groups := CollateActions(txn.InitAction)
	for a := tst.tst.State.TxnActionIdx; a < len(groups); a++ {
		var pa *planaction.PlanAction
		// ************* pre kahuna!  **********
		g := groups[a]

		// don't pull from the incoming queue if we're unable to handle the context.
		for i, v := range g.action {
			if v.Type == "url" {
				if ctx == nil {
					ctx = q.GetUrl()
				}
				if ctx != nil {
					g.action[i].Args["_context"] = ctx
				} else { // if ctx still== nil
					delete(g.action[i].Args, "_context")
					// no url, continuing will just break stuff.
					break
				}
			} else {
				ctx = nil
				delete(g.action[i].Args, "_context")
			}
		}

		// only one action indicates either a malformed satisfygroup (!) or no satisfygroup.
		// assume no.
		if len(g.action) == 1 {
			logger.Debugf("Only one action in group, skipping satisfy test")
			pa = g.action[0]
		} else {
			// we have a choice.
			for _, v := range g.action {
				if v.Type == "url" && ctx == nil {
					// no url received.
					logger.Debugf("No URL received for URL group, skipping.")
					return
				}
				// the only case that err would be set was handled above.
				ac, _ := actions.Satisfy(v.Type, v.Args)
				logger.Tracef("satisfy returned %v", ac)
				if ac == true {
					var m mapstructure.Metadata
					pa = v
					e1, ok := v.Args["on_expected"]
					if ok {
						err = mapstructure.DecodeMetadata(e1, &ex, &m)
						if err != nil {
							logger.Warningf("Could not decode on_expected, defaulting to transaction: %s", err)
						}
						if len(m.Unused) > 0 {
							logger.Warningf("Unused keys in mapstructure decode: %+v", m.Unused)
						}
					}
				}
				logger.Tracef("not satisfied, continuing.")
			}
		}

		if pa == nil {
			logger.Debugf("a satisfygroup %s was declared, but none of the associated actions matched.", g.name)
			logger.Debugf("not executing.")
			tst.processing = false
			return
		} else {
			// Check to see if we got an abort command.
			select {
			case <-kill:
				logger.Warningf("Got command to abort, doing so.")
				for _, v := range tst.bgActions {
					logger.Warningf("aborting action %s", v.GetName())
					v.Abort()
				}
				tst.bgActions = make([]actions.Action, 0)
				tst.processing = false
				return
			default:
				// no abort command received, just keep going.
			}

			// ************* ^^^^^^^^^^^  **********
			logger.Tracef("ex: %+v nex: %+v", ex, nex)
			logger.Tracef("Running action %s", pa.Type)

			// ************* big kahuna!  **********
			logger.Tracef("calling action %s with %v", pa.Type, pa.Args)
			action, res := actions.Execute(pa.Type, pa.Args, tst.tst)
			tst.action = action
			// ************** ^^^^^^^^^ *************
			if res.Err != nil {
				tst.tst.State.States[len(tst.tst.State.States)-1].Status = "errored"
				logger.Warningf("Error running action: %s: %s", pa.Type, res.Err)
				tst.tst.State.Err = res.Err
				tst.processing = false
				return
			}
			if action.CanBackground() {
				logger.Debugf("This action can background, appending to list: %s", action.GetName())
				tst.bgActions = append(tst.bgActions, action)
			}
			if !res.Complete {
				tst.tst.State.States[len(tst.tst.State.States)-1].Status = "waiting"
				logger.Tracef("action %s in process, but not complete", pa.Type)
				tst.processing = false
				return
			} else {
				logger.Tracef("action %s completed", pa.Type)

				if res.Advance {
					logger.Debugf("Advancing to %s", res.NewTxn)
					tst.tst.State.States[len(tst.tst.State.States)-1].Status = "completed"
					err := tst.tst.Advance(res.NewTxn)
					if err != nil {
						logger.Warningf("Error advancing: %s", err)
						tst.tst.State.Err = res.Err
						tst.processing = false
						return
					}
					tst.processing = false
					return
				}
				urlres = &res
				tst.tst.State.TxnActionIdx++
				continue
			}
		}
	}
	// we're at the end of the list of actions - either that action was a url or it was something
	// else.  If it's a url, we want to take actions (expected or unexpected) based upon whether
	// the url succeeded or failed.
	// note:  don't put a url in the expected or nonexpected.  Bad things could happen.  Or at least
	// unexpected.
	if ctx == nil {
		logger.Tracef("Since there was no url, there is no on_expected or on_unexpected.  Stopping.")
		tst.processing = false
		return
	}
	if urlres == nil {
		logger.Tracef("no url result, returning")
		tst.processing = false
		return
	}

	txn, err = tst.tst.GetCurrentTransaction()
	if err != nil {
		logger.Warningf("no current transaction??")
		tst.processing = false
		return
	}
	var e expected.Expected
	if ctx != nil && urlres.Success {
		//e = txn.OnExpected
		e = ex
		tst.tst.State.States[len(tst.tst.State.States)-1].Status = "expected"
	} else {
		//e = txn.OnUnexpected
		e = nex
		tst.tst.State.States[len(tst.tst.State.States)-1].Status = "unexpected"
	}
	code, err := strconv.Atoi(e.ResponseCode)
	if err != nil {
		logger.Warningf("Invalid response code (not int)")
		logger.Debugf("error: %s", err)
		tst.tst.State.Err = err
		tst.processing = false
		return
	}
	ctx.Ctx.Writer.WriteHeader(code)
	
	// Validate response file path to prevent path traversal
	if err := security.ValidatePath(e.Response, ""); err != nil {
		logger.Warningf("Response file path validation failed: %s", err)
		tst.tst.State.Err = err
		tst.processing = false
		return
	}
	
	rawresponsestr, err := os.ReadFile(e.Response)
	if err != nil {
		logger.Warningf("Couldn't read response file %s: %s", e.Response, err)
		tst.tst.State.Err = err
		tst.processing = false
		return
	}

	responsestr := ParseStringTemplate(tst.tst, string(rawresponsestr))

	ctx.Ctx.Writer.Write([]byte(responsestr))
	ctx.Finished <- true

	for _, pa := range e.Action {
		// don't record the action.  That's really only useful for interruptible actions,
		// and I can't think of any reason to do a callback here...
		_, res := actions.Execute(pa.Type, pa.Args, tst.tst)
		if !res.Complete {
			return
		}
		if res.Advance {
			logger.Debugf("advancing to transaction %s", res.NewTxn)
			tst.tst.Advance(res.NewTxn)
			tst.processing = false
			return
		}
	}

	logger.Warningf("expected/unexpected action had no advance")
	tst.tst.State.Err = errors.New("no advance action specified")

}

func FindSatisfaction(n string, g []*SatisfyGroup) *SatisfyGroup {
	for i, v := range g {
		if v.name == n {
			return g[i]
		}
	}
	return nil
}

func CollateActions(a []planaction.PlanAction) []*SatisfyGroup {
	logger := loggo.GetLogger("default")
	s := make([]*SatisfyGroup, 0)

	for i, v := range a {
		// this action is not a part of a SatisfyGroup.  Put it in an empty SatisfyGroup entry and move on.
		if v.SatisfyGroup == "" {
			logger.Tracef("Action v has no satisfy group, adding singleton")
			g := SatisfyGroup{
				name:   "", // this can and should be empty.
				action: []*planaction.PlanAction{&a[i]},
			}
			s = append(s, &g)
			continue
		}
		// this action is part of a SatisfyGroup.
		e := FindSatisfaction(v.SatisfyGroup, s)
		if e == nil {
			logger.Tracef("Action has satisfy group %s, but doesn't exist - adding.", v.SatisfyGroup)
			g := SatisfyGroup{
				name:   v.SatisfyGroup,
				action: []*planaction.PlanAction{&a[i]},
			}
			s = append(s, &g)
			continue
		} else {
			logger.Tracef("Action has satisfy group %s, already exists - adding.", v.SatisfyGroup)
			e.action = append(e.action, &a[i])
		}
	}
	return s
}

func GetPlan() *plan.Plan {
	return tst.tst
}

// forklifted from actions.  If I used it from actions, it would create an import loop.
func ParseStringTemplate(p *plan.Plan, in string) string {

	logger := loggo.GetLogger("default")

	t := time.Now()
	mt, err := t.MarshalText()
	if err != nil {
		logger.Criticalf("Couldn't marshal time: %s", err.Error())
		panic("AIEEEEEEEEE")
	}

	// Anonymous arg template struct.
	at := struct {
		Variables map[string]string
		Bases     map[string]string
		Now       string
	}{
		Bases: p.Bases,
		Now:   string(mt),
	}

	at.Variables = make(map[string]string, 0)

	// Only template the string variables.  It doesn't make much sense
	// for any other type.
	for i, v := range p.State.Variables {
		if reflect.TypeOf(v).String() == "string" {
			at.Variables[i] = v.(string)
		}
	}
	tt := template.New("local")
	tt.Delims("<<", ">>")
	tpl, err := tt.Parse(in)
	if err != nil {
		logger.Warningf("Error parsing template: %s", err.Error())
		return in
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, at)
	if err != nil {
		logger.Warningf("Error executing template: %s.  It's still alive.", err.Error())
		return in
	}
	logger.Tracef("returning: %s", b.String())
	return b.String()
}
