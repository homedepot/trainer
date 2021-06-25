package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"errors"
	"fmt"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"reflect"
)

type Math struct {
	Action
	Args ArgStruct
}

func (w *Math) GetName() string {
	return "math"
}
func (w *Math) Abort() {
	return
}

func (w *Math) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing math action")

	r.Complete = true

	variable, err := w.Args.GetArg("variable", reflect.TypeOf(""), true)
	if err != nil {
		r.Err = fmt.Errorf("couldnt get variable: %w", err)
		return
	}

	valiface, ok := w.Args.Args["value"]
	if !ok {
		r.Err = errors.New("value attribute not specified in math operation")
		return
	}
	valfloat, err := IfaceToFloat(valiface)
	if err != nil {
		r.Err = err
		return
	}

	variface, err := p.State.GetVariable(variable.(string))
	if err != nil {
		r.Err = fmt.Errorf("couldnt get variable %s: %w", variable.(string), err)
		return
	}

	varfloat, err := IfaceToFloat(variface)
	if err != nil {
		r.Err = err
		return
	}

	operation, err := w.Args.GetArg("action", reflect.TypeOf(""), true)
	if err != nil {
		r.Err = fmt.Errorf("couldnt get action arg: %w", err)
		return
	}

	o := MathOps{
		LeftOp:  varfloat,
		RightOp: valfloat,
	}

	result, err := o.Execute(operation.(string))
	if err != nil {
		r.Err = fmt.Errorf("couldnt execute math operation %s: %w", operation.(string), err)
		return
	}

	// if there's already a variable and it's an int, the float64 will be casted back.  Delete it so there
	// no automatic casting.  If it comes in as an int, fine.  It can't go out that way.
	delete(p.State.Variables, variable.(string))
	err = p.State.SetVariable(variable.(string), result)
	if err != nil {
		r.Err = fmt.Errorf("couldnt set variable %s: %w", variable.(string), err)
	} else {
		r.Success = true
	}
	return
}

func (m *Math) SetArgs(i map[string]interface{}) {
	m.Args.Args = i
}

func (m *Math) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Math action: there are no conditions to satisfy")
}

func (m *Math) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (m *Math) CanBackground() bool {
	return false
}

func (m *Math) IsBackgrounded() bool {
	return false
}
