package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"errors"
	"fmt"
	"github.com/homedepot/trainer/security"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"os"
	"reflect"
)

type Match struct {
	Action
	Args ArgStruct
}

func (m *Match) GetName() string {
	return "match"
}
func (m *Match) Abort() {
	return
}

func (m *Match) Execute(p *plan.Plan) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing match action")

	r.Complete = true
	matchfile, err := m.Args.GetArg("match_file", reflect.TypeOf(""), true)
	if err != nil {
		r.Err = err
		return
	}

	matchfiletype, err := m.Args.GetArg("match_file_type", reflect.TypeOf(""), true)

	comp, err := m.LoadFile(matchfile.(string), matchfiletype.(string))
	if err != nil {
		r.Err = err
		return
	}

	responsetype, err := m.Args.GetArg("response_type", reflect.TypeOf(""), true)
	if err != nil {
		r.Err = err
		return
	}

	variable, err := m.Args.GetArg("variable", reflect.TypeOf(""), true)
	if err != nil {
		r.Err = err
		return
	}

	response, err := p.State.GetVariable(variable.(string))
	if err != nil {
		r.Err = err
		return
	}

	respstring, ok := response.(string)
	if !ok {
		r.Err = errors.New("response variable is not a string")
		return
	}
	resp, err := m.LoadString(respstring, responsetype.(string))
	if err != nil {
		r.Err = err
		return
	}

	// Identify if result is a match.
	res := MatchingInterfaces(comp, resp)
	logger.Tracef("MatchingInterfaces() result: %b", res)

	var advanceTxn interface{}
	if res == true {
		// advance to the txn in match_success
		advanceTxn, err = m.Args.GetArg("advance_true", reflect.TypeOf(""), true)
		if err != nil {
			r.Err = err
			return
		}
	} else {
		// advance to the txn in match_failure
		advanceTxn, err = m.Args.GetArg("advance_false", reflect.TypeOf(""), true)
		if err != nil {
			r.Err = err
			return
		}
	}
	logger.Tracef("Executing advance")
	r.Advance = true
	r.NewTxn = advanceTxn.(string)
	/*err = p.Advance(advanceTxn.(string))
	if err != nil {
		r.Err = err
		return
	}*/

	r.Success = res
	return
}

func (m *Match) LoadString(b string, t string) (interface{}, error) {
	logger := loggo.GetLogger("default")

	var f func(string) (interface{}, error)

	switch t {
	case "json":
		logger.Tracef("Loading %s as JSON", b)
		f = LoadJSON
	case "yaml":
		logger.Tracef("Comparing %s as YAML", b)
		f = LoadYAML
	case "string":
		logger.Tracef("Comparing %s as string", b)
		f = LoadString
	default:
		return nil, errors.New(fmt.Sprintf("invalid match file type %s", t))
	}
	i, err := f(b)
	if err != nil {
		return nil, err
	}
	return i, nil
}
func (m *Match) LoadFile(n string, t string) (interface{}, error) {

	logger := loggo.GetLogger("default")

	var f func(string) (interface{}, error)

	// Validate match file path to prevent path traversal
	if err := security.ValidatePath(n, ""); err != nil {
		logger.Criticalf("match file path validation failed: %s", err.Error())
		return nil, err
	}

	// Capture match file.
	body, err := os.ReadFile(n)
	if err != nil {
		logger.Criticalf("match file read failed: %s", err.Error())
		return nil, err
	}

	switch t {
	case "json":
		logger.Tracef("Loading %s as JSON", n)
		f = LoadJSON
	case "yaml":
		logger.Tracef("Loading %s as YAML", n)
		f = LoadYAML
	case "string":
		logger.Tracef("Loading %s as string", n)
		f = LoadString
	default:
		return nil, errors.New(fmt.Sprintf("invalid match file type %s", t))
	}
	i, err := f(string(body))
	if err != nil {
		return nil, err
	}
	return i, nil
}

// MatchingInterfaces identifies if interface
// types and interface contents match.
func MatchingInterfaces(mFile interface{}, candidate interface{}) bool {
	logger := loggo.GetLogger("default")

	// Capture the type of our input, and
	// compare it to our pared interface type.
	switch t := mFile.(type) {
	case map[string]interface{}:
		if reflect.TypeOf(candidate).String() != "map[string]interface {}" {
			logger.Warningf("Mismatched types: expecting map[string]interface {}, got %s", reflect.TypeOf(candidate).String())
			return false
		}
		c := candidate.(map[string]interface{})
		for k, v := range mFile.(map[string]interface{}) {
			_, ok := c[k]
			if !ok {
				logger.Errorf("required key missing in pared struct")
				return false
			}
			if !MatchingInterfaces(v, c[k]) {
				return false
			}
		}
		return true
	case []interface{}:
		if reflect.TypeOf(candidate).String() != "[]interface {}" {
			logger.Warningf("Mismatched types: expecting []interface {}, got %s", reflect.TypeOf(candidate).String())
			return false
		}
		c := candidate.([]interface{})
		for k, v := range mFile.([]interface{}) {
			if !MatchingInterfaces(v, c[k]) {
				return false
			}
		}
		return true
	case string:
		if reflect.TypeOf(candidate).Name() != "string" {
			logger.Warningf("Mismatched types: expecting string, got %s", reflect.TypeOf(candidate).Name())
			return false
		}
		res := mFile.(string) == candidate.(string)
		if res == false {
			logger.Tracef("string didn't match: file: '%s' response: '%s'", mFile.(string), candidate.(string))
		}
		return res
	case bool:
		if reflect.TypeOf(candidate).Name() != "bool" {
			logger.Warningf("Mismatched types: expecting bool, got %s", reflect.TypeOf(candidate).Name())
			return false
		}
		res := mFile.(bool) == candidate.(bool)
		if res == false {
			logger.Tracef("bool didn't match: %v %v", mFile.(bool), candidate.(bool))
		}
		return res
	case int:
		if reflect.TypeOf(candidate).Name() != "int" {
			logger.Warningf("Mismatched types: expecting int, got %s", reflect.TypeOf(candidate).Name())
			return false
		}
		return mFile.(int) == candidate.(int)
	case float64:
		if reflect.TypeOf(candidate).Name() != "float64" {
			logger.Warningf("Mismatched types: expecting float64, got %s", reflect.TypeOf(candidate).Name())
			return false
		}
		return mFile.(float64) == mFile.(float64)
	default:
		logger.Errorf("invalid case %s", t)
		return false
	}
}

func (m *Match) SetArgs(i map[string]interface{}) {
	m.Args.Args = i
}

func (m *Match) Satisfy() (bool, error) {
	return false, errors.New("cannot use satisfy_group for Match action: there are no conditions to satisfy")
}

func (m *Match) GetContext() (*context.Context, *context.CancelFunc) {
	return nil, nil
}

func (m *Match) CanBackground() bool {
	return false
}

func (m *Match) IsBackgrounded() bool {
	return false
}
