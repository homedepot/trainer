package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/homedepot/trainer/cli"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"github.com/stretchr/testify/assert"
)

var Config *config.Config
var Username = "blah"
var Password = "blah"
var o cli.Options

func TestMain(m *testing.M) {
	// The tempdir is created so MongoDB has a location to store its files.
	// Contents are wiped once the server stops
	err := os.Setenv("APIAUTHUSERNAME", Username)
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	err = os.Setenv("APIAUTHPASSWORD", Password)
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	err = os.Setenv("TESTMODE", "true")
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	err = os.Setenv("LOGLEVEL", "TRACE")
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	err = os.Setenv("CONFIGFILE", "../data/config.yml")
	if err != nil {
		log.Printf("Couldn't setenv: %s", err.Error())
		return
	}
	o = *cli.NewOptions()
	o.Parse([]string{})

	Config = LoadConfig(o)
	loglevel, berr := loggo.ParseLevel(o.LogLevel)
	if berr == false {
		log.Printf("Couldn't parse loglevel %s", o.LogLevel)
		return
	}
	loggo.GetLogger("default").SetLogLevel(loglevel)

	// Run the test suite
	m.Run()
}

func LoadConfig(opt cli.Options) *config.Config {
	logger := loggo.GetLogger("default")

	absPath, err := filepath.Abs(opt.ConfigFile)
	o.ConfigFile = absPath
	if err = GoToConfigDir(opt); err != nil {
		logger.Criticalf("Unable to find configuration directory: %s", err)
		os.Exit(1)
	}

	c, err := config.Load(absPath, opt.TestMode, opt.TestURL, opt.Bases)
	if err != nil {
		logger.Criticalf("Couldn't load config: %s", err.Error())
		os.Exit(1)
	}
	return c
}

// FindConfigDir finds the directory for the configuration
// based on the file path and configuration file path value.
func FindConfigDir(opt cli.Options) string {
	logger := loggo.GetLogger("default")
	logger.Debugf("CONFIGFILE = %q", opt.ConfigFile)
	configDir := filepath.Dir(opt.ConfigFile)
	if configDir == "" {
		logger.Criticalf("Generated config directory path is an empty string")
		os.Exit(1)
	}
	logger.Debugf("ConfigDir = %q", configDir)
	return configDir
}

// GoToConfigDir executes a change directory given
// the path received from call to FindConfigDir().
func GoToConfigDir(opt cli.Options) error {
	logger := loggo.GetLogger("default")
	configDir := FindConfigDir(opt)
	if err := os.Chdir(configDir); err != nil {
		return err
	}
	logger.Debugf("Changed working directory to %q", configDir)
	return nil
}

func TestCompareInterfaces(t *testing.T) {

	injson := `{"1": "2"}`
	paredjson := `{"1": "2", "2": "3"}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.True(t, result, "Not true")
	t.Logf("The computer says yes: %t", result)
}

func TestCompareInterfacesUnequal(t *testing.T) {

	injson := `{"1": "2"}`
	paredjson := `{"1": "3", "2": "3"}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.False(t, result, "True")
}

func TestCompareInterfacesNonexistent(t *testing.T) {

	injson := `{"1": "2"}`
	paredjson := `{"2": "3"}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.False(t, result, "True")
}

func TestCompareInterfacesComplex(t *testing.T) {

	injson := `{"1": {"2": "3"}}`
	paredjson := `{"1": {"2": "3", "4": "5"}}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.True(t, result, "False")
}

func TestCompareInterfacesBool(t *testing.T) {

	injson := `{"1": {"2": true}}`
	paredjson := `{"1": {"2": true, "4": "5"}}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.True(t, result, "False")
}

func TestCompareInterfacesBoolUnequal(t *testing.T) {

	injson := `{"1": {"2": true}}`
	paredjson := `{"1": {"2": false, "4": "5"}}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.False(t, result, "True")
}

func TestCompareInterfacesMismatch(t *testing.T) {

	injson := `{"1": {"2": true}}`
	paredjson := `{"1": {"2": "mismatch", "4": "5"}}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.False(t, result, "True")
}

func TestCompareInterfacesComplexWithArray(t *testing.T) {

	injson := `{"1": [ {"2": "3"}, {"4": "5"}]}`
	paredjson := `{"1": [ {"2": "3"}, {"4": "5"}, {"6": "7"}]}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.True(t, result, "False")
}

func TestCompareInterfacesComplexWithMismatchedArray(t *testing.T) {

	injson := `{"1": [ {"2": "3"}, {"4": "5"}]}`
	paredjson := `{"1": {"2": "3"}}`

	jint, err := LoadJSON(injson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	pint, err := LoadJSON(paredjson)
	if err != nil {
		t.Errorf("It's an error!: %s", err.Error())
		t.Fail()
		return
	}

	result := MatchingInterfaces(jint, pint)
	assert.False(t, result, "True")
}

func TestExecute(t *testing.T) {
	pl, err := Config.FindPlan("basic_test")
	if err != nil {
		t.Errorf("couldn't load config")
		return
	}
	type args struct {
		t string
		a map[string]interface{}
		p *plan.Plan
	}
	tests := []struct {
		name string
		args args
		want ExecuteResult
	}{
		{
			name: "execute test function",
			args: args{
				t: "test",
				a: make(map[string]interface{}),
				p: pl,
			},
			want: ExecuteResult{
				Complete: true,
				Success:  true,
				Err:      nil,
			},
		},
	}
	SetLogger("DEBUG")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := Execute(tt.args.t, tt.args.a, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
