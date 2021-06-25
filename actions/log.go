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

type Log struct {
	Action
	Args         ArgStruct
	DefaultLevel string
}

func (l *Log) GetName() string {
	return "log"
}
func (l *Log) Abort() {
	return
}

func (l *Log) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing log action")

	r.Complete = true

	value, err := l.Args.GetArg("value", reflect.TypeOf(""), true)
	if err != nil {
		logger.Warningf("couldn't get log value: %s", err)
		r.Err = err
		return
	}
	loglevel, err := l.Args.GetArg("log_level", reflect.TypeOf(""), false)
	if err != nil {
		logger.Warningf("couldn't get log level: %s", err)
		r.Err = err
		return
	}
	if loglevel == nil {
		loglevel = "l.DefaultLevel"
	}
	v := value.(string)
	ll := loglevel.(string)
	lv, ok := loggo.ParseLevel(ll)
	if ok == false {
		r.Err = errors.New("invalid log level")
		return
	}
	logger.Logf(lv, v)
	r.Success = true
	return

}

func (l *Log) SetArgs(i map[string]interface{}) {
	l.Args.Args = i
}

func (l *Log) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Log action: there are no conditions to satisfy")
}

func (l *Log) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (l *Log) CanBackground() bool {
	return false
}

func (l *Log) IsBackgrounded() bool {
	return false
}
