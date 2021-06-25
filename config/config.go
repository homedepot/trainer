package config

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"errors"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/homedepot/trainer/structs/state"
	"github.com/homedepot/trainer/structs/transaction"
	"github.com/juju/loggo"
	"github.com/mohae/deepcopy"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config struct for housing configuration data,
// specifically a default plan, set of provided plans,
// and a set of base URLs for testing. The test URL will
// be added to Bases, whether the Config data contains
// additional base URLs in its Bases map or not.
type Config struct {
	DefaultPlan  string            `yaml:"default_plan" json:"default_plan"`
	Plans        []plan.Plan       `yaml:"plan" json:"plan"`
	PlanIncludes []string          `yaml:"planinclude" json:"planinclude"`
	Bases        map[string]string `yaml:"bases" json:"bases"`
}

// NewConfig creates a new configuration given
// a config file, testing mode, and test URL.
// tm (testing mode) iteratively included, but not yet applied.
func NewConfig(cfile string, tm bool, tURL string, bases map[string]string) (*Config, error) {

	config := Config{}

	// Read configuration file and then
	// marshal the YAML into our []byte.
	b, err := ioutil.ReadFile(cfile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	// If a test URL present, add
	// it to Bases map of config.
	if config.Bases == nil {
		config.Bases = make(map[string]string, 0)
	}
	if tURL != "" || tm {
		config.Bases["testurl"] = tURL
	}
	if bases != nil && len(bases) > 0 {
		for k, v := range bases {
			config.Bases[k] = v
		}
	}

	err = config.LoadPlans()
	if err != nil {
		return nil, err
	}

	err = config.LoadVars()
	if err != nil {
		return nil, err
	}

	err = config.ValidateConfig()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate config validates Config.
// See Issue #21
// TODO resolve issues (based upon schema file of some kind)
// TODO: - Douglas (we do not want those ifs at all, so we need to do it right and validate the yaml)
func (c *Config) ValidateConfig() error {
	for i, _ := range c.Plans {
		for j, _ := range c.Plans[i].Txn {

			if c.Plans[i].Bases == nil {
				c.Plans[i].Bases = make(map[string]string, 0)
			}

			if len(c.Plans[i].Bases) == 0 {
				c.Plans[i].Bases = c.Bases
			} else {
				c.Plans[i].Bases["testurl"] = c.Bases["testurl"]
			}

			/*if c.Plans[i].Txn[j].URL == "" {
				c.Plans[i].Txn[j].Standalone = true
				continue
			}
			*/

			// Temporarily pass over any configurations that
			// are missing required data. Only standalone
			// transactions should have empty data here. If
			// transaction is not standalone, but is missing
			// response data, we're currently just skipping
			// them. Iteratively, we will add validation to
			// the yaml to ensure this is handled eloquently.
			if c.Plans[i].Txn[j].OnExpected.Response == "" {
				continue
			}
			if c.Plans[i].Txn[j].OnUnexpected.Response == "" {
				continue
			}

			c.Plans[i].Txn[j].OnExpected.Expected = true
			c.Plans[i].Txn[j].OnUnexpected.Expected = false

		}
	}
	return nil
}

// FindPlan locates a plan for configuration in our Config.
func (c *Config) FindPlan(name string) (*plan.Plan, error) {
	for i, _ := range c.Plans {
		if c.Plans[i].Name == name {
			return &c.Plans[i], nil
		}
	}
	return nil, errors.New("failed to locate plan " + name)

}

// LoadConfig takes a config file and sets
// the pipeline configuration to it. This
// dictates the behavior of the testing.
func Load(cfile string, testmode bool, testurl string, bases map[string]string) (*Config, error) {
	logger := loggo.GetLogger("default")

	c, err := NewConfig(cfile, testmode, testurl, bases)
	if err != nil {
		logger.Criticalf("Couldn't load config: %s", err.Error())
		return nil, err
	}
	return c, err
}

// GetFirstTransaction returns a Config's first plan transaction.
func (c *Config) GetFirstTransaction(plan string) (*transaction.Transaction, error) {
	p, err := c.FindPlan(plan)
	if err != nil {
		return nil, err
	}
	return p.GetFirstTxn()
}

func (c *Config) LoadPlans() error {
	logger := loggo.GetLogger("default")
	logger.Debugf("Loading plans...")
	if c.PlanIncludes == nil {
		return nil
	}

	for i, _ := range c.PlanIncludes {
		logger.Debugf("loading plan file %s", c.PlanIncludes[i])
		err := c.LoadPlanFile(c.PlanIncludes[i])
		if err != nil {
			logger.Warningf("Couldn't load plan file %s: %s", c.PlanIncludes[i], err)
			return err
		}
	}

	for k, _ := range c.Plans {
		logger.Debugf("Loading transactions")
		err := c.Plans[k].LoadTransactions()
		if err != nil {
			logger.Warningf("Couldn't load transactions: %s", err)
		}
		// generate any on-the-fly transaction stuffs needed (like URL actions)
		c.Plans[k].InitializeTransactions()
	}
	return nil
}

func (c *Config) LoadPlanFile(name string) error {
	logger := loggo.GetLogger("default")
	str, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	if len(str) == 0 {
		logger.Warningf("Empty plan file: %s", str)
		return nil
	}
	i := plan.Plan{}
	err = yaml.Unmarshal(str, &i)
	if err != nil {
		logger.Warningf("Error reading plan file: %s", err.Error())
		return err
	}

	logger.Debugf("plan: %+v", i)
	c.Plans = append(c.Plans, i)
	return nil
}

func (c *Config) LoadVars() error {
	for i, _ := range c.Plans {

		if c.Plans[i].State == nil {
			txn, err := c.Plans[i].GetFirstTxn()
			if err != nil {
				return err
			}
			c.Plans[i].State = state.NewState(txn.Name)
		}

		if c.Plans[i].DefaultVars == nil {
			c.Plans[i].DefaultVars = make(map[string]interface{}, 0)
		}

		if c.Plans[i].ExtVarFile != "" {
			str, err := ioutil.ReadFile(c.Plans[i].ExtVarFile)
			if err != nil {
				return err
			}

			var m map[string]string
			err = yaml.Unmarshal(str, &m)
			if err != nil {
				return err
			}

			for k, v := range m {
				c.Plans[i].DefaultVars[k] = v
			}
		}
		res := deepcopy.Copy(c.Plans[i].DefaultVars)
		c.Plans[i].State.Variables = res.(map[string]interface{})
	}
	return nil
}
