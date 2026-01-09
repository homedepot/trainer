package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/homedepot/trainer/security"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"strconv"
	"text/template"
	"time"
)

func IfaceToFloat(i interface{}) (float64, error) {
	logger := loggo.GetLogger("default")
	var err error
	var retval float64
	t := reflect.TypeOf(i).String()
	if t == "string" {
		retval, err = strconv.ParseFloat(i.(string), 64)
	} else if t == "float64" {
		retval = i.(float64)
		err = nil
	} else if t == "float32" {
		retval = float64(i.(float32))
	} else if t == "int" {
		retval = float64(i.(int))
	} else {
		return 0, errors.New("unknown type for math comparison " + t)
	}
	if err != nil {
		logger.Warningf("Couldn't parse float")
		return 0, err
	}
	return retval, nil
}

// LoadJSON unmarshals JSON into an
// interface and returns the interface.
func LoadJSON(in string) (interface{}, error) {
	var i interface{}
	err := json.Unmarshal([]byte(in), &i)
	if err != nil {
		return nil, err
	}
	return i, err
}

// LoadYAML unmarshals YAML/YML into an
// interface and returns the interface.
func LoadYAML(in string) (interface{}, error) {
	var i interface{}
	err := yaml.Unmarshal([]byte(in), &i)
	if err != nil {
		return nil, err
	}
	return i, err
}

// LoadString just returns a string.
func LoadString(in string) (interface{}, error) {
	return in, nil
}

func ParseTemplate(p *plan.Plan, args map[string]interface{}) (map[string]interface{}, error) {

	logger := loggo.GetLogger("default")

	t := time.Now()
	mt, err := t.MarshalText()
	if err != nil {
		logger.Criticalf("Couldn't marshal time: %s", err.Error())
		panic("AIEEEEEEEEE")
	}

	// Anonymous arg template struct.
	at := struct {
		Variables map[string]interface{}
		Bases     map[string]string
		Now       string
	}{
		Bases: p.Bases,
		Now:   string(mt),
	}

	at.Variables = p.State.Variables
	out := make(map[string]interface{}, 0)

	// Only template the string variables.  It doesn't make much sense
	// for any other type.
	tt := template.New("local")
	tt.Delims("<<", ">>")
	for j, w := range args {
		str, ok := w.(string)
		if !ok {
			// passthrough, it's not a string
			out[j] = args[j]
			continue
		}
		tpl, err := tt.Parse(str)
		if err != nil {
			logger.Warningf("Error parsing template: %s", err.Error())
			return nil, err
		}
		var b bytes.Buffer
		err = tpl.Execute(&b, at)
		if err != nil {
			logger.Warningf("Error executing template: %s.  It's still alive.", err.Error())
			return nil, err
		}
		out[j] = b.String()
	}
	logger.Tracef("returning: %v", out)
	return out, nil
}

// EqualStrings returns a boolean value
// for string comparison, as well as error.
func EqualStrings(input string, expected string) (bool, interface{}, error) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Entering the incomparably boring EqualStrings")
	if input == "" {
		return false, nil, errors.New("string value for input empty")
	}
	if expected == "" {
		return false, nil, errors.New("string value for expected empty")
	}
	return input == expected, nil, nil
}

// EqualJSON returns a boolean value
// for JSON comparison, as well as error.
func EqualJSON(input string, expected string) (bool, interface{}, error) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Entering CompareJSON")

	// Interfaces for comparison.
	var a, b interface{}

	err := json.Unmarshal([]byte(input), &a)
	if err != nil {
		return false, nil, err
	}

	// we should refactor this at some point so that this doesn't happen every time.
	err = json.Unmarshal([]byte(expected), &b)
	if err != nil {
		return false, nil, err
	}

	logger.Tracef("a: %+v", a)
	logger.Tracef("b: %+v", b)
	//res := reflect.DeepEqual(a, b)
	res := MatchingInterfaces(b, a)
	logger.Tracef("res: %+b", res)
	return res, a, nil
}

// EqualYAML returns a boolean value
// for YAML comparison, as well as error.
func EqualYAML(input string, expected string) (bool, interface{}, error) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Entering EqualYAML")

	// Interfaces for comparison.
	var a, b interface{}

	err := yaml.Unmarshal([]byte(input), &a)
	if err != nil {
		return false, nil, err
	}

	// we should refactor this at some point so that this doesn't happen every time.
	err = yaml.Unmarshal([]byte(expected), &b)
	if err != nil {
		return false, nil, err
	}

	return reflect.DeepEqual(a, b), a, nil
}

// ExecuteComparison compares return body string and
// OnExpected value for a State's transaction.
func ExecuteComparison(p *plan.Plan, data string, path string, strbody string, f func(string, string) (bool, interface{}, error)) (interface{}, error) {

	logger := loggo.GetLogger("default")

	// Validate data file path to prevent path traversal
	if err := security.ValidatePath(data, ""); err != nil {
		logger.Criticalf("Data file path validation failed: %s", err.Error())
		return nil, err
	}

	// Get expected data.
	expected, err := os.ReadFile(data)
	if err != nil {
		// Could not read file. Fail the test.
		// TODO should we add some Slack posting functionality here? (Douglas)
		logger.Criticalf("Couldn't read transaction data: %s: %s", data, err.Error())
		logger.Criticalf("Terminating test run for %s", path)
		return nil, err
	}

	// Execute comparison.
	tequ := ParseStringTemplate(p, string(expected))
	logger.Tracef("tequ: %s, strbody: %s", tequ, strbody)
	equal, i, err := f(strbody, tequ)
	if err != nil {
		logger.Tracef("Error comparing input for %s: %s", path, err.Error())
		return nil, err
	}
	if !equal {
		logger.Tracef("Comparison failed ---> not equal!")
		return i, errors.New("comparison failed ---> not equal")
	}

	return i, nil
}

func ParseStringTemplate(p *plan.Plan, in string) string {

	logger := loggo.GetLogger("default")

	t := time.Now()
	mt, err := t.MarshalText()
	if err != nil {
		logger.Criticalf("Couldn't marshal time: %s", err.Error())
		panic("AIEEEEEEEEE")
	}

	// Anonymous arg template struct.
	at := struct {
		Variables map[string]string
		Bases     map[string]string
		Now       string
	}{
		Bases: p.Bases,
		Now:   string(mt),
	}

	at.Variables = make(map[string]string, 0)

	// Only template the string variables.  It doesn't make much sense
	// for any other type.
	for i, v := range p.State.Variables {
		if reflect.TypeOf(v).String() == "string" {
			at.Variables[i] = v.(string)
		}
	}
	tt := template.New("local")
	tt.Delims("<<", ">>")
	tpl, err := tt.Parse(in)
	if err != nil {
		logger.Warningf("Error parsing template: %s", err.Error())
		return in
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, at)
	if err != nil {
		logger.Warningf("Error executing template: %s.  It's still alive.", err.Error())
		return in
	}
	logger.Tracef("returning: %s", b.String())
	return b.String()
}
