package expected

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/homedepot/trainer/structs/planaction"
)

// TODO comment this
type Expected struct {
	Response     string                  `yaml:"response" json:"response"`
	ResponseCode string                  `yaml:"response_code" json:"response_code" mapstructure:"response_code"`
	ResponseType string                  `yaml:"response_contenttype" json:"respoonse_contenttype" mapstructure:"response_code"`
	Action       []planaction.PlanAction `yaml:"action" json:"action"`
	Expected     bool                    `yaml:"-" json:"-"`
}
