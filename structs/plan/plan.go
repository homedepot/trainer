package plan

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"bytes"
	"errors"
	"github.com/homedepot/trainer/structs/state"
	"github.com/homedepot/trainer/structs/transaction"
	"github.com/juju/loggo"
	"github.com/mohae/deepcopy"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"text/template"
	"time"
)

// TODO comment this
type Plan struct {
	Name             string                    `yaml:"name" json:"name"`
	DefaultVars      map[string]interface{}    `yaml:"variables" json:"variables"`
	ExtVarFile       string                    `yaml:"externalvars" json:"externalvars"`
	Txn              []transaction.Transaction `yaml:"txn" json:"txn"`
	CurrentTimer     int                       `yaml:"-" json:"-"` // This gets incremented every second
	ExpectedTimer    int                       `yaml:"-" json:"-"`
	Bases            map[string]string         `yaml:"bases" json:"bases"`
	StartTransaction string                    `yaml:"start_transaction" json:"start_transaction"`
	TxnIncludes      []TxnInclude              `yaml:"txninclude" json:"txninclude"`
	StopVar          string                    `yaml:"stop_var" json:"stop_var"`
	State            *state.State              `yaml:"state" json:"state"`
}

type TxnInclude struct {
	File      string
	Variables map[string]interface{}
}

// GetFirstTxn returns a plan's first transaction.
func (p *Plan) GetFirstTxn() (*transaction.Transaction, error) {
	if p == nil {
		return nil, errors.New("plan is nil, would panic")
	}
	if len(p.Txn) == 0 {
		return nil, errors.New("no transactions")
	}
	if p.StartTransaction != "" {
		return p.FindTransaction(p.StartTransaction)
	}
	return &p.Txn[0], nil
}

func (p *Plan) LoadTransactions() error {
	logger := loggo.GetLogger("default")
	logger.Debugf("Loading transactions...")
	if p.TxnIncludes == nil {
		return nil
	}

	for i, _ := range p.TxnIncludes {
		logger.Debugf("loading transaction file %s", p.TxnIncludes[i])
		err := p.LoadTxnFile(p.TxnIncludes[i])
		if err != nil {
			logger.Warningf("Couldn't load transaction file %s: %s", p.TxnIncludes[i], err)
			return err
		}
	}
	return nil
}

func (p *Plan) InitializeTransactions() {
	logger := loggo.GetLogger("default")
	logger.Debugf("Loading transactions...")
	for i, _ := range p.Txn {
		p.Txn[i].CreateUrlAction()
	}
	return
}

func (p *Plan) LoadTxnFile(txn TxnInclude) error {
	logger := loggo.GetLogger("default")
	str, err := ioutil.ReadFile(txn.File)
	if err != nil {
		return err
	}

	at := struct {
		Variables map[string]string
	}{}

	at.Variables = make(map[string]string, 0)
	for i, v := range txn.Variables {
		if reflect.TypeOf(v).String() == "string" {
			at.Variables[i] = v.(string)
		}
	}
	if len(at.Variables) > 0 {
		logger.Debugf("Parsing template...")
		tt := template.New("local1")
		tt.Delims("[[", "]]")
		tpl, err := tt.Parse(string(str))
		if err != nil {
			logger.Warningf("Error parsing template: %s", err.Error())
			return err
		}
		var b bytes.Buffer
		err = tpl.Execute(&b, at)
		if err != nil {
			logger.Warningf("Error executing template: %s.  It's still alive.", err.Error())
			return err
		}
		str = b.Bytes()
	}

	if len(str) == 0 {
		logger.Warningf("Empty transaction file: %s", str)
		return nil
	}
	i := transaction.Transaction{}
	err = yaml.Unmarshal(str, &i)
	if err != nil {
		logger.Warningf("Error reading plan file: %s", err.Error())
		return err
	}
	logger.Debugf("transaction: %+v", i)
	p.Txn = append(p.Txn, i)

	// set the variables in the main state with the txninclude defaults.
	return nil
}

func (p *Plan) FindTransaction(name string) (*transaction.Transaction, error) {
	for i, _ := range p.Txn {
		if p.Txn[i].Name == name {
			return &p.Txn[i], nil
		}
	}
	return nil, errors.New("no such transaction")
}

func (p *Plan) Reset() error {
	txn, err := p.GetFirstTxn()
	if err != nil {
		return err
	}
	p.CurrentTimer = 0

	p.State = state.NewState(txn.Name)

	err = p.State.Reset(txn.Name)
	if err != nil {
		return err
	}

	res := deepcopy.Copy(p.DefaultVars)
	p.State.Variables = res.(map[string]interface{})

	return nil
}

func (p *Plan) Advance(name string) error {
	logger := loggo.GetLogger("default")
	t, err := p.FindTransaction(name)
	if err != nil {
		return err
	}

	logger.Infof("Entering transaction %s", t.Name)
	p.State.NewState(t.Name)

	return nil
}

func (p *Plan) GetCurrentTransaction() (*transaction.Transaction, error) {
	t, err := p.FindTransaction(p.State.Transaction)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (p *Plan) ParseTemplate(in string) string {

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
	logger.Debugf("returning: %s", b.String())
	return b.String()
}
