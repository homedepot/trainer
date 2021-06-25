package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"errors"
	"github.com/homedepot/trainer/structs/plan"
)

type Callback struct {
	Action
	Args ArgStruct
	ctx  context.Context
	cf   context.CancelFunc
}

func (c *Callback) GetName() string {
	return "callback"
}
func (c *Callback) Abort() {
	return
}

func (c *Callback) Execute(p *plan.Plan) (r ExecuteResult) {
	ctx := context.Background()
	cancelctx, cf := context.WithCancel(ctx)
	c.ctx = cancelctx
	c.cf = cf
	return DoCallback(c.Args, p, ctx)
}

func (c *Callback) SetArgs(i map[string]interface{}) {
	c.Args.Args = i
}

func (c *Callback) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Callback action: there are no conditions to satisfy")
}

func (c *Callback) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (c *Callback) CanBackground() bool {
	return false
}

func (c *Callback) IsBackgrounded() bool {
	return false
}
