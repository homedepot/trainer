package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"reflect"
)

type ExecuteResult struct {
	Err          error
	Complete     bool
	Success      bool
	Advance      bool
	NewTxn       string
	RegisterURL  string
	RegisterUUID uuid.UUID
}

type Actions struct {
	Actions map[string]Action
}

// implementation note:
// Satisfy is implemented across all actions, but it only makes sense in some.
// Actions that do not have an internal state do not make sense to advance, as there
// is nothing that will allow us to determine whether a condition is satisfied.
type Action interface {
	Execute(*plan.Plan) ExecuteResult
	Abort()
	GetName() string
	SetArgs(map[string]interface{})
	Satisfy() (bool, error)
	GetContext() (*context.Context, *context.CancelFunc)
	CanBackground() bool
	IsBackgrounded() bool
}

func SetLogger(t string) {
	l, ok := loggo.ParseLevel(t)
	if !ok {
		// only should be used for testing anyway
		panic(fmt.Sprintf("invalid loglevel: %s", t))
	}
	logger := loggo.GetLogger("default")
	logger.SetLogLevel(l)
}

type ArgStruct struct {
	Args map[string]interface{}
}

type QueueContext struct {
	Ctx      *gin.Context
	Finished chan bool
}

var ActionsArr = map[string]Action{
	"advance":     &Advance{},
	"callback":    &Callback{},
	"cbsplit":     &CbSplit{},
	"cbfinish":    &CbFinish{},
	"conditional": &Conditional{},
	"log":         &Log{},
	"match":       &Match{},
	"math":        &Math{},
	"set":         &Set{},
	"wait":        &Wait{},
	"test":        &Test{},
	"url":         &URL{},
}

func NewActions() *Actions {
	a := &Actions{}
	a.Actions = ActionsArr
	return a
}

func (a *ArgStruct) SetArg(n string, i interface{}) error {
	if a.Args == nil {
		a.Args = make(map[string]interface{}, 0)
	}
	i1, ok := a.Args[n]
	if ok {
		if reflect.TypeOf(i) != reflect.TypeOf(i1) {
			return errors.New("differing interface types")
		}
	}
	a.Args[n] = i
	return nil
}

func (a *ArgStruct) GetArg(n string, t reflect.Type, required bool) (interface{}, error) {
	if a.Args == nil {
		a.Args = make(map[string]interface{}, 0)
	}

	if a.Args[n] == nil && !required {
		return nil, nil
	}

	if reflect.TypeOf(a.Args[n]) != t {
		return nil, errors.New(fmt.Sprintf("unexpected type %s", reflect.TypeOf(a.Args[n])))
	}
	arg, ok := a.Args[n]
	if !ok {
		if required {
			return nil, errors.New(fmt.Sprintf("argument %s not found", n))
		} else {
			return nil, nil
		}
	}
	return arg, nil
}

// DoAction executes the respective action requested, executing
// callbacks, advancing, waiting, or resetting as needed.
func Execute(t string, a map[string]interface{}, p *plan.Plan) (Action, ExecuteResult) {

	logger := loggo.GetLogger("default")
	logger.Tracef("starting execute execute: action %s", t)
	as := NewActions()
	action := as.Actions[t]
	action.SetArgs(a)
	return action, action.Execute(p)
}

func Satisfy(t string, a map[string]interface{}) (bool, error) {
	logger := loggo.GetLogger("default")
	logger.Tracef("starting execute execute: action %s", t)
	as := NewActions()
	action := as.Actions[t]
	action.SetArgs(a)
	return action.Satisfy()
}
