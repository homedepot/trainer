package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"errors"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"reflect"
)

type Advance struct {
	Action
	Args ArgStruct
}

func (a *Advance) GetName() string {
	return "advance"
}
func (a *Advance) Abort() {
	// noop, there's nothing to abort, it advances or doesn't.
	return
}
func (a *Advance) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing advance action")

	r.Complete = true // no matter what happens, it will always complete
	txn, err := a.Args.GetArg("txn", reflect.TypeOf(""), true)
	if err != nil {
		r.Err = err
		return
	}
	// don't try if there's no such transaction
	_, err = p.FindTransaction(txn.(string))
	if err != nil {
		r.Err = err
		return
	}

	r.Advance = true
	r.NewTxn = txn.(string)
	r.Success = true
	return
}

func (a *Advance) SetArgs(i map[string]interface{}) {
	a.Args.Args = i
}

func (a *Advance) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Advance action: there are no conditions to satisfy")
}

func (a *Advance) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (a *Advance) CanBackground() bool {
	return false
}

func (a *Advance) IsBackgrounded() bool {
	return false
}
