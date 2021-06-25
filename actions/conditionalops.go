package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"errors"
	"fmt"
	"reflect"
)

type ConditionalOps struct {
	LeftOp   interface{}
	RightOp  interface{}
	cmpFuncs map[string]func() (bool, error)
}

func NewConditionalOps() *ConditionalOps {
	c := &ConditionalOps{}
	if c.cmpFuncs == nil {
		c.cmpFuncs = map[string]func() (bool, error){
			"eq": c.CompareEq,
			"ne": c.CompareNe,
			"lt": c.CompareLt,
			"gt": c.CompareGt,
			"le": c.CompareLe,
			"ge": c.CompareGe,
		}
	}
	return c
}
func (c *ConditionalOps) Compare(op string) (bool, error) {

	_, ok := c.cmpFuncs[op]
	if !ok {
		return false, errors.New(fmt.Sprintf("compare: invalid operator %s", op))
	}
	l, ok1 := c.LeftOp.(string)
	r, ok2 := c.RightOp.(string)

	if ok1 && ok2 {
		return c.cmpFuncs[op]()
	} else if ok1 != ok2 {
		return false, errors.New(fmt.Sprintf("invalid conversion: %s to %s",
			reflect.TypeOf(l).String(),
			reflect.TypeOf(r).String()))
	}

	lb, ok1 := c.LeftOp.(bool)
	rb, ok1 := c.LeftOp.(bool)

	if ok1 && ok2 {
		return c.cmpFuncs[op]()
	} else if ok1 != ok2 {
		return false, errors.New(fmt.Sprintf("invalid conversion: %s to %s",
			reflect.TypeOf(lb).String(),
			reflect.TypeOf(rb).String()))
	}

	success, err := c.cmpFuncs[op]()
	if err != nil {
		return false, err
	} else {
		return success, nil
	}
}

func (c *ConditionalOps) ConvertToFloat64() (float64, float64, error) {

	var l, r float64
	if reflect.TypeOf(c.LeftOp).String() == "int" {
		l = float64(c.LeftOp.(int))
	} else if reflect.TypeOf(c.LeftOp).String() == "float32" {
		l = float64(c.LeftOp.(float32))
		r = float64(c.RightOp.(float32))
	} else if reflect.TypeOf(c.LeftOp).String() == "float64" {
		l = c.LeftOp.(float64)
	} else {
		return 0, 0, errors.New(fmt.Sprintf("cannot convert %s to type float64", reflect.TypeOf(c.LeftOp).String()))
	}

	if reflect.TypeOf(c.RightOp).String() == "int" {
		r = float64(c.RightOp.(int))
	} else if reflect.TypeOf(c.RightOp).String() == "float32" {
		r = float64(c.RightOp.(float32))
	} else if reflect.TypeOf(c.RightOp).String() == "float64" {
		r = c.RightOp.(float64)
	} else {
		return 0, 0, errors.New(fmt.Sprintf("cannot convert %s to type float64", reflect.TypeOf(c.RightOp).String()))
	}

	return l, r, nil
}

func (c *ConditionalOps) CompareEq() (bool, error) {
	_, ok1 := c.RightOp.(bool)
	_, ok2 := c.LeftOp.(bool)
	if ok1 != ok2 {
		return false, errors.New("one boolean and one something else")
	}
	if ok1 && ok2 {
		l := c.LeftOp.(bool)
		r := c.RightOp.(bool)
		return l == r, nil
	}
	_, ok1 = c.LeftOp.(string)
	_, ok2 = c.RightOp.(string)
	if ok1 && ok2 {
		l := c.LeftOp.(string)
		r := c.RightOp.(string)
		return l == r, nil
	}
	l, r, err := c.ConvertToFloat64()
	if err != nil {
		return false, err
	}
	return l == r, nil
}

func (c *ConditionalOps) CompareNe() (bool, error) {
	_, ok1 := c.RightOp.(bool)
	_, ok2 := c.LeftOp.(bool)
	if ok1 != ok2 {
		return false, errors.New("one boolean and one something else")
	}
	if ok1 && ok2 {
		l := c.LeftOp.(bool)
		r := c.RightOp.(bool)
		return l != r, nil
	}
	_, ok1 = c.LeftOp.(string)
	_, ok2 = c.RightOp.(string)
	if ok1 && ok2 {
		l := c.LeftOp.(string)
		r := c.RightOp.(string)
		return l != r, nil
	}
	l, r, err := c.ConvertToFloat64()
	if err != nil {
		return false, err
	}
	return l != r, nil
}

func (c *ConditionalOps) CompareLt() (bool, error) {
	_, ok1 := c.RightOp.(bool)
	_, ok2 := c.LeftOp.(bool)
	if ok1 && ok2 {
		return false, errors.New(fmt.Sprintf("lt: cannot order booleans"))
	}
	_, ok1 = c.LeftOp.(string)
	_, ok2 = c.RightOp.(string)
	if ok1 && ok2 {
		l := c.LeftOp.(string)
		r := c.RightOp.(string)
		return l < r, nil
	}

	l, r, err := c.ConvertToFloat64()
	if err != nil {
		return false, err
	}
	return l < r, nil
}

func (c *ConditionalOps) CompareLe() (bool, error) {
	_, ok1 := c.RightOp.(bool)
	_, ok2 := c.LeftOp.(bool)
	if ok1 && ok2 {
		return false, errors.New(fmt.Sprintf("le: cannot order booleans"))
	}
	_, ok1 = c.LeftOp.(string)
	_, ok2 = c.RightOp.(string)
	if ok1 && ok2 {
		l := c.LeftOp.(string)
		r := c.RightOp.(string)
		return l <= r, nil
	}

	l, r, err := c.ConvertToFloat64()
	if err != nil {
		return false, err
	}
	return l <= r, nil
}

func (c *ConditionalOps) CompareGe() (bool, error) {
	_, ok1 := c.RightOp.(bool)
	_, ok2 := c.LeftOp.(bool)
	if ok1 && ok2 {
		return false, errors.New(fmt.Sprintf("ge: cannot order booleans"))
	}
	_, ok1 = c.LeftOp.(string)
	_, ok2 = c.RightOp.(string)
	if ok1 && ok2 {
		l := c.LeftOp.(string)
		r := c.RightOp.(string)
		return l >= r, nil
	}
	l, r, err := c.ConvertToFloat64()
	if err != nil {
		return false, err
	}
	return l >= r, nil
}

func (c *ConditionalOps) CompareGt() (bool, error) {
	_, ok1 := c.RightOp.(bool)
	_, ok2 := c.LeftOp.(bool)
	if ok1 && ok2 {
		return false, errors.New(fmt.Sprintf("gt: cannot order booleans"))
	}
	_, ok1 = c.LeftOp.(string)
	_, ok2 = c.RightOp.(string)
	if ok1 && ok2 {
		l := c.LeftOp.(string)
		r := c.RightOp.(string)
		return l > r, nil
	}

	l, r, err := c.ConvertToFloat64()
	if err != nil {
		return false, err
	}
	return l > r, nil
	return false, errors.New(fmt.Sprintf("gt: incompatible types: %s and %s",
		reflect.TypeOf(c.LeftOp).String(),
		reflect.TypeOf(c.RightOp).String()))
}
