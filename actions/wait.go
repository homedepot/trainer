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
	"strconv"
	"time"
)

type Wait struct {
	Action
	Args         ArgStruct
	Reset        bool
	CurrentTimer int
}

func (w *Wait) GetName() string {
	return "wait"
}
func (w *Wait) Abort() {
	w.Reset = true
	return
}

func (w *Wait) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing wait action")

	if p.State.AbortRunningAction {
		r.Complete = true
		r.Err = errors.New("aborted")
		return
	}
	if p.State.WaitActionStartTime.IsZero() {
		p.State.WaitActionStartTime = time.Now()
	}
	starttime := p.State.WaitActionStartTime

	duration, err := w.Args.GetArg("duration", reflect.TypeOf(0), true)
	if err != nil {
		logger.Warningf("%s", err)
		r.Err = err
		return
	}
	durint := duration.(int)

	// Manage sleep and reset cycle.

	// if reset is true, something called abort.
	seconds := time.Since(starttime)
	durstr := strconv.Itoa(durint) + "s"
	durat, err := time.ParseDuration(durstr)
	if err != nil {
		r.Err = err
		return
	}
	if seconds >= durat {
		logger.Debugf("Wait completed!")
		p.State.WaitActionStartTime = time.Time{}
		r.Success = true
		r.Complete = true
		return
	} else {
		logger.Tracef("Wait not completed (%v %v)", seconds, durint)
		return
	}
}

func (w *Wait) SetArgs(i map[string]interface{}) {
	w.Args.Args = i
}

func (w *Wait) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Wait action: there are no conditions to satisfy")
}

func (w *Wait) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (w *Wait) CanBackground() bool {
	return false
}

func (w *Wait) IsBackgrounded() bool {
	return false
}
