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

type CbFinish struct {
	Action
	Args ArgStruct
}

func (c *CbFinish) GetName() string {
	return "cbfinish"
}
func (c *CbFinish) Abort() {
	return
}

func (c *CbFinish) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing cb_finish action")
	if !currcb.inprogress {
		logger.Warningf("no split callback to finish")
		return ExecuteResult{
			Success: false,
		}
	}
	logger.Debugf("Waiting for completion of background process...")
	aborted := <-currcb.aborted
	logger.Debugf("Background process completed.  Aborted = %b", aborted)
	var out *ExecuteResult
	if *aborted {
		out = &ExecuteResult{
			Complete: true,
			Success:  false,
		}
		logger.Warningf("split callback aborted, failing")
	} else {
		logger.Debugf("Not aborted, waiting for output")
		out = <-currcb.output
		logger.Debugf("output received.")
	}
	currcb.inprogress = false
	logger.Infof("cb_finish: Execute returns %+v", out)
	return *out
}

func (c *CbFinish) SetArgs(i map[string]interface{}) {
	c.Args.Args = i
}

func (c *CbFinish) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Callback action: there are no conditions to satisfy")
}

func (c *CbFinish) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (c *CbFinish) CanBackground() bool {
	return false
}

func (c *CbFinish) IsBackgrounded() bool {
	return false
}
