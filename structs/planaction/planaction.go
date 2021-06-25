package planaction

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

// TODO comment this
type PlanAction struct {
	Type         string                 `yaml:"type" json:"type"`
	SatisfyGroup string                 `yaml:"satisfy_group" json:"satisfy_group"`
	Args         map[string]interface{} `yaml:"args" json:"args"`
}
