package transaction

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/homedepot/trainer/structs/expected"
	"github.com/homedepot/trainer/structs/planaction"
)

// TODO comment this
type Transaction struct {
	Name          string                  `yaml:"name" json:"name"`
	URL           string                  `yaml:"url" json:"url"`
	Data          string                  `yaml:"data" json:"data"`
	Datatype      string                  `yaml:"datatype" json:"datatype"`
	OnExpected    expected.Expected       `yaml:"on_expected" json:"on_expected"`
	OnUnexpected  expected.Expected       `yaml:"on_unexpected" json:"on_unexpected"`
	InitAction    []planaction.PlanAction `yaml:"init_action" json:"init_action"`
	Standalone    bool                    `yaml:"-" json:"-"` // only consists of actions, no web responses
	Running       bool                    `yaml:"-" json:"-"` // set when the transaction is running - only if standalone
	SaveBody      string                  `yaml:"save_body" json:"save_body"`
	SaveBodyAsMap string                  `yaml:"save_body_as_map" json:"save_body_as_map"`
}

func (t *Transaction) CreateUrlAction() {
	if t.URL != "" {
		a := planaction.PlanAction{}
		args := map[string]interface{}{
			"url": t.URL,
		}
		// initializing in the initializer adds the key whether or not the string is empty.
		// that is incorrect behavior.
		if t.SaveBody != "" {
			args["save_body"] = t.SaveBody
		}
		if t.SaveBodyAsMap != "" {
			args["save_body_as_map"] = t.SaveBodyAsMap
		}
		if t.Data != "" {
			args["data"] = t.Data
		}
		if t.Datatype != "" {
			args["datatype"] = t.Datatype
		}
		a.Type = "url"
		a.Args = args
		t.InitAction = append(t.InitAction, a)
		t.Standalone = false
	} else {
		t.Standalone = true
	}
}
