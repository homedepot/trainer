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
)

type CbSplit struct {
	Action
	Args ArgStruct
	ctx  context.Context
	cf   context.CancelFunc
}

func (c *CbSplit) GetName() string {
	return "cbsplit"
}
func (c *CbSplit) Abort() {
	currcb.inprogress = false
	c.cf()
	return
}

func (c *CbSplit) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing cb_split action")
	if currcb.inprogress {
		// not entirely sure what to do here.
		// panicking would kill the whole thing.
		// but we don't know what to advance to on failure.  I don't think...  TODO
		return ExecuteResult{
			Err:      errors.New("already a split callback in progress"),
			Complete: true,
			Success:  false,
		}
	}
	currcb.inprogress = true
	output := make(chan *ExecuteResult)
	aborted := make(chan *bool)
	currcb.output = output
	currcb.aborted = aborted
	ctx := context.Background()
	cancelctx, cancelfunc := context.WithCancel(ctx)
	c.ctx = cancelctx
	c.cf = cancelfunc

	go func() {
		r := DoCallback(c.Args, p, cancelctx)
		logger.Debugf("Callback finished, determining whether it was cancelled...")
		var aborted bool
		select {
		case <-cancelctx.Done():
			// the callback was cancelled
			logger.Debugf("Callback cancelled")
			aborted = true
			currcb.aborted <- &aborted
		default:
			logger.Debugf("Callback not cancelled")
			aborted = false
			currcb.aborted <- &aborted
		}
		logger.Debugf("Callback was completed, finishing callback")
		output <- &r
	}()
	logger.Tracef("finished cb_split execute")
	// aborted may be misnamed for this purpose, but it's an accurate description of what this determins.  It just
	// has a side effect of also telling us when the callback is done.
	return ExecuteResult{
		Complete: true,
		Success:  true,
	}
}

func (c *CbSplit) SetArgs(i map[string]interface{}) {
	c.Args.Args = i
}

func (c *CbSplit) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Callback action: there are no conditions to satisfy")
}

func (c *CbSplit) GetContext() (*context.Context, *context.CancelFunc) {
	return &c.ctx, &c.cf
}

func (c *CbSplit) CanBackground() bool {
	return true
}

func (c *CbSplit) IsBackgrounded() bool {
	return currcb.inprogress
}
