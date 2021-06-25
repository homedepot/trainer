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

type Set struct {
	Action
	Args ArgStruct
}

func (s *Set) GetName() string {
	return "wait"
}
func (s *Set) Abort() {
	return
}

func (s *Set) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing set action")

	r.Complete = true

	variable, err := s.Args.GetArg("variable", reflect.TypeOf(""), true)
	if err != nil {
		logger.Warningf("couldn't get variable: %s", err)
		r.Err = err
		return
	}

	source, ok1 := s.Args.Args["source"]
	value, ok2 := s.Args.Args["value"]

	if ok1 && source.(string) != "" {
		v, err := p.State.GetVariable(source.(string))
		if err != nil {
			logger.Warningf("couldn't get variable %s: err", source.(string), err)
			r.Err = err
			return
		}
		if reflect.TypeOf(v).String() == "string" {
			err = p.State.SetVariable(variable.(string), v.(string))
			if err != nil {
				r.Err = err
				return
			}
		} else if reflect.TypeOf(v).String() == "int" {
			err = p.State.SetVariable(variable.(string), v.(int))
			if err != nil {
				r.Err = err
			}
		} else if reflect.TypeOf(v).String() == "float64" {
			err = p.State.SetVariable(variable.(string), v.(float64))
			if err != nil {
				r.Err = err
				return
			}
		} else if reflect.TypeOf(v).String() == "float32" {
			err = p.State.SetVariable(variable.(string), v.(float32))
			if err != nil {
				r.Err = err
				return
			}
		} else if reflect.TypeOf(v).String() == "bool" {
			err = p.State.SetVariable(variable.(string), v.(bool))
			if err != nil {
				r.Err = err
				return
			}
		}
	} else if ok2 {
		err := p.State.SetVariable(variable.(string), value)
		if err != nil {
			r.Err = err
			return
		}
	} else {
		r.Err = errors.New("value or source not defined")
		return
	}
	r.Success = true
	return
}

func (s *Set) SetArgs(i map[string]interface{}) {
	s.Args.Args = i
}

func (s *Set) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Set action: there are no conditions to satisfy")
}

func (s *Set) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (s *Set) CanBackground() bool {
	return false
}

func (s *Set) IsBackgrounded() bool {
	return false
}
