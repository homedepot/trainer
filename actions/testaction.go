package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
)

type Test struct {
	Action
	Args  ArgStruct
	Reset bool
}

func (t *Test) GetName() string {
	return "test"
}

func (t *Test) Abort() {
	t.Reset = true
	return
}

// This is a no-op solely for testing.
func (t *Test) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing test action")

	r = ExecuteResult{
		Complete: true,
		Success:  true,
		Err:      nil,
	}
	return
}

func (t *Test) SetArgs(i map[string]interface{}) {
	t.Args.Args = i
}

func (a *Test) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}
