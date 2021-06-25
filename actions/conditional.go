package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"errors"
	"fmt"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"reflect"
)

type Conditional struct {
	Action
	Args ArgStruct
}

func (c *Conditional) GetName() string {
	return "conditional"
}
func (c *Conditional) Abort() {
	return
}

func (c *Conditional) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing conditional action")

	r.Complete = true

	rawterm, ok := c.Args.Args["term"]
	if !ok {
		r.Err = errors.New("term is nil")
		return
	}

	t := make(map[interface{}]interface{})
	term, ok := rawterm.(map[interface{}]interface{})
	if ok {
		t = term
	} else {
		w, ok := rawterm.(map[string]interface{})
		if ok {
			for k, v := range w {
				t[k] = v
			}
		} else {
			r.Err = errors.New("term isn't a useful map")
			return
		}
	}

	fmt.Printf("t: %v", t)
	variable, ok := t["variable"].(string)
	if !ok {
		r.Err = errors.New("no variable specified on left side of term")
		return
	}
	conditional, ok := t["conditional"].(string)
	if !ok {
		r.Err = errors.New("no conditional specified in term")
		return
	}

	var leftop interface{}
	var rightop interface{}
	var err error

	cvar, ok1 := t["conditional_var"].(string)
	cval, ok2 := t["conditional_value"]

	if !ok1 && !ok2 {
		r.Err = errors.New("must specify one of conditional_var or conditional_value")
		return
	} else if ok1 {
		rightop, err = p.State.GetVariable(cvar)
		if err != nil {
			logger.Warningf("Specified undeclared variable for conditional variable: %s", err)
			r.Err = err
			return
		}
	} else {
		rightop = cval
	}

	logger.Tracef("Variables: %s", p.State.Variables)
	leftop, err = p.State.GetVariable(variable)
	if err != nil {
		logger.Warningf("Specified undeclared variable for conditional operation: %s", err)
		r.Err = err
		return
	}

	if leftop == nil {
		logger.Warningf("Trying to compare a nil variable: %s", variable)
		r.Err = errors.New("nil variable")
		return
	}

	cond := NewConditionalOps()

	cond.LeftOp = leftop
	cond.RightOp = rightop

	result, err := cond.Compare(conditional)
	if err != nil {
		r.Err = err
		return
	}

	logger.Tracef("Conditional:  leftop: %v (%s) rightop: %v (%s) operation: %s result: %v", cond.LeftOp, reflect.TypeOf(cond.LeftOp).String(), cond.RightOp, reflect.TypeOf(cond.RightOp).String(), conditional, result)
	var advanceTxn interface{}
	if result == true {
		// advance to the txn in match_success
		advanceTxn, err = c.Args.GetArg("advance_true", reflect.TypeOf(""), true)
		if err != nil {
			r.Err = errors.New("advance_true not set")
			return
		}
	} else {
		// advance to the txn in match_failure
		advanceTxn, err = c.Args.GetArg("advance_false", reflect.TypeOf(""), true)
		if err != nil {
			r.Err = errors.New("advance_false not set")
			return
		}
	}
	logger.Tracef("Executing advance")
	r.Advance = true
	r.NewTxn = advanceTxn.(string)

	r.Success = result
	return
}

func (c *Conditional) SetArgs(i map[string]interface{}) {
	c.Args.Args = i
}

func (c *Conditional) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Conditional action: there are no conditions to satisfy")
}

func (c *Conditional) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (c *Conditional) CanBackground() bool {
	return false
}

func (c *Conditional) IsBackgrounded() bool {
	return false
}
