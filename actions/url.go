package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"io"
	"reflect"
)

type URL struct {
	Action
	Args       ArgStruct
	Registered bool
	uuid       uuid.UUID
}

func (u *URL) GetName() string {
	return "url"
}
func (u *URL) Abort() {
	return
}

func (u *URL) URLMatches(ctx *QueueContext) (bool, error) {
	logger := loggo.GetLogger("default")
	logger.Tracef("checking existing context...")
	url, err := u.Args.GetArg("url", reflect.TypeOf(""), true)
	if err != nil {
		return false, err
	}
	logger.Tracef("comparing %s and %s", ctx.Ctx.Request.URL.Path, url.(string))
	if ctx.Ctx.Request.URL.Path != url.(string) {
		logger.Warningf("unexpected URL %s", ctx.Ctx.Request.URL.Path)
		return false, nil
	}
	return true, nil

}
func (u *URL) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing url action")

	context, ok := u.Args.Args["_context"]
	if !ok {
		logger.Debugf("no URL received (empty context)")
		return
	}
	ctx := context.(*QueueContext)
	res, err := u.URLMatches(ctx)
	if err != nil {
		r.Complete = true
		r.Err = err
		return
	}
	if res == false {
		r.Complete = true
		return
	}
	body, err := io.ReadAll(ctx.Ctx.Request.Body)
	if err != nil {
		logger.Criticalf("Couldn't read request body, Failing test. (%s)", err.Error())
		r.Err = err
		return
	}
	strbody := string(body)
	logger.Tracef("Got body %s", strbody)
	savebody, ok := u.Args.Args["save_body"]
	if ok {
		sb, ok := savebody.(string)
		if ok {
			if sb != "" {
				logger.Tracef("saving body into %s", sb)
				p.State.Variables[sb] = strbody
			}
		} else {
			r.Err = errors.New("savebody var isn't a string")
			return
		}
	}

	// Read the expected data from the specified file if one
	// is specified.  If not, assume the body doesn't matter.
	var i interface{}
	var failed bool
	data, ok := u.Args.Args["data"]
	if ok {
		dt, ok := data.(string)
		if ok {
			var f func(string, string) (bool, interface{}, error)

			dtype := ""
			datatype, ok := u.Args.Args["datatype"]
			if ok {
				dstr, ok := datatype.(string)
				if ok && dstr != "" {
					dtype = dstr
				}
			}
			switch d := dtype; d {
			case "json":
				logger.Tracef("Comparing as JSON")
				f = EqualJSON
			case "yaml":
				logger.Tracef("Comparing as YAML")
				f = EqualYAML
			default:
				logger.Tracef("Nothing compares to you --Sinead")
				f = EqualStrings
			}

			// Compare the returned body with the
			// transaction's OnExpected value.
			i, err = ExecuteComparison(p, dt, ctx.Ctx.Request.URL.Path, strbody, f)
			if err != nil {
				logger.Warningf("%v failed ExecuteComparison()", ctx.Ctx.Request.URL.Path)
				failed = true
			} else {
				logger.Tracef("Comparison succeeded")
			}
		}
	}

	sbam, ok := u.Args.Args["save_body_as_map"]
	if ok {
		sbamstr, ok := sbam.(string)
		if ok && sbamstr != "" {
			logger.Tracef("saving body into map %s", sbamstr)
			p.State.Variables[sbamstr] = i
		}
	}

	if failed {
		r.Complete = true
		return
	}
	logger.Tracef("finished...")
	r.Complete = true
	r.Success = true

	return
}

func (u *URL) SetArgs(i map[string]interface{}) {
	u.Args.Args = i
}

func (u *URL) Satisfy() (bool, error) {
	// return true if the url received matches the url expected
	logger := loggo.GetLogger("default")
	logger.Tracef("satisfy: checking existing context...")
	context, ok := u.Args.Args["_context"]
	if !ok {
		logger.Debugf("no URL received (empty context)")
		return false, nil
	}
	ctx := context.(*QueueContext)
	return u.URLMatches(ctx)
}

func (u *URL) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (u *URL) CanBackground() bool {
	return false
}

func (u *URL) IsBackgrounded() bool {
	return false
}
