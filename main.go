package main

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"fmt"
	"github.com/homedepot/trainer/cli"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/handler"
	"github.com/homedepot/trainer/metrics"
	"github.com/homedepot/trainer/router"
	"github.com/juju/loggo"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var o cli.Options

// Version is updated in build (see /ci/script/build.sh).
// For now, GitCommit and BuildStamp are included, as
// they may yet be used in a later version of build.sh.
var (
	Version     = "0.0.1"
	GitCommit   = "HEAD"
	BuildStamp  = "UNKNOWN"
	FullVersion = Version + "+" + GitCommit + "-" + BuildStamp
)

// init sets up versioning and Prometheus.
func init() {
	metrics.Register()
}

// PrintVars prints captured Kingpin vars.
func PrintVars(o cli.Options) {
	// Using log because we haven't yet set up the logger.
	log.Printf("Starting application with following vars:")
	log.Printf("     Version: %s", Version)
	log.Printf("     Listening on: %s:%s", o.APIListenHost, o.APIListenPort)
	log.Printf("		Config File: %s", o.ConfigFile)
	log.Printf("     LogLevel: %s", o.LogLevel)
	if o.TestMode {
		log.Printf("     TEST MODE ENABLED")
	} else {
		log.Printf("     PRODUCTION MODE")
	}
}

func main() {

	o = *cli.NewOptions()
	o.Parse(os.Args[1:])
	o.App.Version(FullVersion)
	logLevel, err := loggo.ParseLevel(o.LogLevel)
	if err == false {
		log.Printf("failed to parse loglevel %s", o.LogLevel)
		return
	}
	PrintVars(o)
	loggo.GetLogger("default").SetLogLevel(logLevel)
	logger := loggo.GetLogger("default")
	logger.Infof("starting application...")

	// Load configuration and set initial plan for testing.
	c := LoadConfig(o)
	//SetInitialPlan(c)

	// Start testing- managed by ticker and action runner.
	//StartBackgroundTicker(c)
	//InitializeActionRunner(c)
	addr := fmt.Sprintf("%s:%s", o.APIListenHost, o.APIListenPort)
	logger.Infof("listening on " + addr)
	var wg sync.WaitGroup
	h := &handler.Handler{}
	wg.Add(1)
	router.StartRouter(o, c, &wg, h)
	h.Start()
	wg.Wait()
	h.Stop()
}

// LoadConfig takes a config file and sets
// the pipeline configuration to it. This
// dictates the behavior of the testing.
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
