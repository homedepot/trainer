package state

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/juju/loggo"
)

// TODO comment this
type StateEntry struct {
	TxnName string `yaml:"txn_name" json:"txn_name"`
	Status  string `yaml:"status" json:"status"`
}

// TODO comment this
type State struct {
	Transaction         string
	States              []StateEntry
	RunnerKillSwitch    bool
	Variables           map[string]interface{}
	TxnActionIdx        int       // index of completed actions
	TxnActionsCompleted bool      // whether the initactions are entirely completed
	Err                 error     // put any errors here, also blocks any further progress
	WaitActionStartTime time.Time // for waits
	AbortRunningAction  bool
}

func NewState(t string) *State {
	s := &State{}
	s.Transaction = t
	s.Variables = make(map[string]interface{})
	return s
	/*	plan, err := cfg.FindPlan(p)
		if err != nil {
			return nil, fmt.Errorf("error getting plan %s: %w", p, err)
		}
		firsttxn, err := cfg.GetFirstTransaction(p)
		if err != nil {
			return nil, fmt.Errorf( "error getting first txn from %s: %w", p, err)
		}

		s.Plan = plan
		s.Transaction = firsttxn
		return s, nil*/

}

// Reset interrupts current testing actions
// to return state to initial (init) state.
func (s *State) Reset(n string) error {
	s.States = make([]StateEntry, 1)
	s.States[0] = StateEntry{
		TxnName: n,
		Status:  "waiting",
	}
	s.TxnActionIdx = 0
	s.TxnActionsCompleted = false

	return nil
}

// NewState .... TODO update this comment
func (s *State) NewState(name string) {
	s.Transaction = name
	entry := StateEntry{
		TxnName: name,
		Status:  "pending",
	}
	s.States = append(s.States, entry)
	s.TxnActionIdx = 0
	s.TxnActionsCompleted = false
}

func (s *State) GetVariable(varname string) (interface{}, error) {
	logger := loggo.GetLogger("default")
	logger.Debugf("Getting variable: %s", varname)
	arr := ParseString(varname)
	logger.Debugf("Parsed array: %s", arr)
	i, err := s.GetVariableRecursive(nil, arr, nil)
	return i, err
}

func (s *State) GetVariableRecursive(i interface{}, arr []*VariableEntity, value interface{}) (interface{}, error) {
	if i == nil {
		i = s.Variables
	}

	if (arr == nil || len(arr) == 0) && value == nil {
		return nil, errors.New("get recursive var: nil or empty var array")
	}

	// set operation and one variable name left
	if len(arr) == 1 && value != nil {
		v := arr[0]
		if v.Name == "" {
			if reflect.TypeOf(i).String() == "[]interface {}" {
				ri := i.([]interface{})
				n := arr[0].Index
				ri[n] = value
				return nil, nil
			} else {
				return nil, errors.New("not an interface")
			}
		}
		if reflect.TypeOf(i).String() == "map[string]interface {}" {
			ri := i.(map[string]interface{})
			n := arr[0].Name
			ri[n] = value
			return nil, nil
		} else {
			return nil, errors.New("nothing to update")
		}
	}
	v := arr[0]
	arr = arr[1:]

	if v.Name == "" {
		if reflect.TypeOf(i).String() == "[]interface {}" {
			ri := i.([]interface{})
			return s.DoArray(ri, arr, v.Index, value)
		} else {
			return nil, errors.New("not an interface")
		}
	}
	if reflect.TypeOf(i).String() == "map[string]interface {}" {
		ri := i.(map[string]interface{})
		return s.DoMapVariable(ri, arr, v.Name, value)
	} else {
		return s.DoVariable(i, v.Name, value)
	}
}

func (s *State) DoVariable(m interface{}, key string, value interface{}) (interface{}, error) {
	logger := loggo.GetLogger("default")

	if m == nil {
		return nil, errors.New("dovariable value of key is nil")
	}
	mtype := reflect.TypeOf(m).String()
	valuetype := ""
	if value != nil {
		valuetype = reflect.TypeOf(value).String()
	}
	if mtype == "string" {
		logger.Debugf("string")
		if value != nil {
			if valuetype != "string" {
				return nil, errors.New("incompatible type " + valuetype + " to string")
			}
			m = value.(string)
			return "", nil
		} else {
			return m.(string), nil
		}
	} else if mtype == "int" {
		logger.Debugf("int")
		if value != nil {
			if valuetype == "int" {
				m = value.(int)
			} else if valuetype == "float32" {
				m = int(value.(float32))
			} else if valuetype == "float64" {
				m = int(value.(float64))
			}
			return "", nil
		} else {
			return m.(int), nil
		}
	} else if mtype == "float64" {
		logger.Debugf("float64")
		if value != nil {
			if valuetype == "int" {
				m = float64(value.(int))
			} else if valuetype == "float32" {
				m = float64(value.(float32))
			} else if valuetype == "float64" {
				m = value.(float64)
			}
			m = value.(float64)
			return "", nil
		} else {
			return m.(float64), nil
		}
	} else if mtype == "float32" {
		logger.Debugf("float32")
		if value != nil {
			if valuetype == "int" {
				m = float32(value.(int))
			} else if valuetype == "float32" {
				m = value.(float32)
			} else if valuetype == "float64" {
				m = float32(value.(float64))
			}
			return "", nil
		} else {
			return m.(float32), nil
		}
	} else if mtype == "bool" {
		logger.Debugf("bool")
		if value != nil {
			if valuetype != "bool" {
				return nil, errors.New("incompatible type " + valuetype + " to bool")
			}
			m = value.(bool)
			return "", nil
		} else {
			return m.(bool), nil
		}
	} else {
		logger.Warningf("Unknown type %s", mtype)
		return nil, errors.New("unknown type " + mtype)
	}
}

func (s *State) DoMapVariable(m map[string]interface{}, argarr []*VariableEntity, key string, value interface{}) (interface{}, error) {
	logger := loggo.GetLogger("default")

	if m[key] == nil && value != nil {
		logger.Debugf("variable %s does not exist, creating", key)
		if len(argarr) >= 1 {
			logger.Debugf("creating as map[string]interface{}")
			m[key] = make(map[string]interface{}, 0)
		} else {
			m[key] = value
			return nil, nil
		}
	} else if m[key] == nil && value == nil {
		return nil, fmt.Errorf("variable %s does not exist and nothing to set it to", key)
	}
	mtype := reflect.TypeOf(m[key]).String()
	valuetype := ""
	if value != nil {
		valuetype = reflect.TypeOf(value).String()
	}
	if mtype == "map[string]interface {}" || mtype == "[]interface {}" {
		logger.Debugf("%s[%v] (argarr: %+v)", mtype, key, argarr)
		res, err := s.GetVariableRecursive(m[key], argarr, value)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else if reflect.TypeOf(m[key]).String() == "string" {
		logger.Debugf("string")
		if value != nil {
			if valuetype != "string" {
				return nil, errors.New("incompatible type " + valuetype + " to string")
			}
			m[key] = value.(string)
			return "", nil
		} else {
			return m[key].(string), nil
		}
	} else if reflect.TypeOf(m[key]).String() == "int" {
		logger.Debugf("int")
		if value != nil {
			if valuetype == "int" {
				m[key] = value.(int)
			} else if valuetype == "float32" {
				m[key] = int(value.(float32))
			} else if valuetype == "float64" {
				m[key] = int(value.(float64))
			}
			return "", nil
		} else {
			return m[key].(int), nil
		}
	} else if reflect.TypeOf(m[key]).String() == "float64" {
		logger.Debugf("float64")
		if value != nil {
			if valuetype == "int" {
				m[key] = float64(value.(int))
			} else if valuetype == "float32" {
				m[key] = float64(value.(float32))
			} else if valuetype == "float64" {
				m[key] = value.(float64)
			}
			return "", nil
		} else {
			return m[key].(float64), nil
		}
	} else if reflect.TypeOf(m[key]).String() == "float32" {
		logger.Debugf("float32")
		if value != nil {
			if valuetype == "int" {
				m[key] = float32(value.(int))
			} else if valuetype == "float32" {
				m[key] = value.(float32)
			} else if valuetype == "float64" {
				m[key] = float32(value.(float64))
			}
			return "", nil
		} else {
			return m[key].(float32), nil
		}
	} else if reflect.TypeOf(m[key]).String() == "bool" {
		logger.Debugf("bool")
		if value != nil {
			if valuetype != "bool" {
				return nil, errors.New("incompatible type " + valuetype + " to bool")
			}
			m[key] = value.(bool)
			return "", nil
		} else {
			return m[key].(bool), nil
		}
	} else {
		logger.Warningf("Unknown type %s", reflect.TypeOf(m[key]).String())
		return nil, errors.New("unknown type " + reflect.TypeOf(m[key]).String())
	}
}

func (s *State) DoArray(m []interface{}, argarr []*VariableEntity, idx int, value interface{}) (interface{}, error) {
	logger := loggo.GetLogger("default")

	if m == nil {
		// we're refusing to create one because we have no idea how large to make it
		// I wish I was that psychic.
		return nil, errors.New("array does not exist and refusing to create one")
	}
	if idx > len(m) {
		return nil, errors.New("index " + fmt.Sprint(len(m)) + " out of range")
	}
	if m[idx] == nil {
		return nil, errors.New("domapvariable value of key " + fmt.Sprint(idx) + " is nil")
	}
	mtype := reflect.TypeOf(m[idx]).String()
	valuetype := ""
	if value != nil {
		valuetype = reflect.TypeOf(value).String()
	}
	if mtype == "map[string]interface {}" || mtype == "[]interface {}" {
		logger.Debugf("%s[%v] (argarr: %+v)", mtype, idx, argarr)
		res, err := s.GetVariableRecursive(m[idx], argarr, value)
		if err != nil {
			return nil, err
		}
		return res, nil
	} else if mtype == "string" {
		logger.Debugf("string")
		if value != nil {
			if valuetype != "string" {
				return nil, errors.New("incompatible type " + valuetype + " to string")
			}
			m[idx] = value.(string)
			return "", nil
		} else {
			return m[idx].(string), nil
		}
	} else if reflect.TypeOf(m[idx]).String() == "int" {
		logger.Debugf("int")
		if value != nil {
			if value != nil {
				if valuetype == "int" {
					m[idx] = value.(int)
				} else if valuetype == "float32" {
					m[idx] = int(value.(float32))
				} else if valuetype == "float64" {
					m[idx] = int(value.(float64))
				}
			}
			return "", nil
		} else {
			return m[idx].(int), nil
		}
	} else if reflect.TypeOf(m[idx]).String() == "float64" {
		logger.Debugf("float64")
		if value != nil {
			if valuetype == "int" {
				m[idx] = float64(value.(int))
			} else if valuetype == "float32" {
				m[idx] = float64(value.(float32))
			} else if valuetype == "float64" {
				m[idx] = value.(float64)
			}
			return "", nil
		} else {
			return m[idx].(float64), nil
		}
	} else if reflect.TypeOf(m[idx]).String() == "float32" {
		logger.Debugf("float32")
		if value != nil {
			if valuetype == "int" {
				m[idx] = float32(value.(int))
			} else if valuetype == "float32" {
				m[idx] = value.(float32)
			} else if valuetype == "float64" {
				m[idx] = float32(value.(float64))
			}
			return "", nil
		} else {
			return m[idx].(float32), nil
		}
	} else if reflect.TypeOf(m[idx]).String() == "bool" {
		logger.Debugf("bool")
		if value != nil {
			if reflect.TypeOf(m[idx]).String() != "bool" {
				return nil, errors.New("incompatible type " + mtype + " to bool")
			}
			m[idx] = value.(bool)
			return "", nil
		} else {
			return m[idx].(bool), nil
		}
	} else {
		logger.Warningf("Unknown type %s", mtype)
		return nil, errors.New("unknown type " + mtype)
	}
}

func (s *State) SetVariable(varname string, value interface{}) error {
	logger := loggo.GetLogger("default")
	arr := ParseString(varname)
	logger.Debugf("set value: %s", value)
	for i, v := range arr {
		logger.Debugf("  arr[%v] %+v", i, v)
	}
	_, err := s.GetVariableRecursive(nil, arr, value)
	return err
}

type VariableEntity struct {
	Name  string
	Index int
}

func ParseString(in string) []*VariableEntity {
	arr := strings.FieldsFunc(in, func(r rune) bool {
		if r == '.' {
			return true
		}
		if r == '[' {
			return true
		}
		if r == ']' {
			return true
		}
		return false
	})

	var entities []*VariableEntity
	for _, v := range arr {
		e := VariableEntity{}
		result, err := strconv.Atoi(v)
		if err != nil {
			e.Name = v
		} else {
			e.Index = result
		}
		entities = append(entities, &e)
	}

	return entities
}
