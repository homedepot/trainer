package cli

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type Options struct {
	APIListenPort      string
	APIListenHost      string
	APIAuthUsername    string
	APIAuthPass        string
	MongoConnectionURL string
	MongoDBCollection  string
	LogLevel           string
	TestMode           bool
	TestURL            string
	Bases              map[string]string
	ConfigFile         string
	App                *kingpin.Application
}

func NewOptions() *Options {
	o := Options{}
	return &o
}

func (o *Options) Parse(args []string) {
	app := kingpin.New("tentacool", "CRUD storage for snowkemon ecosystem").Author("THD Engineering")

	APIListenPort := app.Flag("apiport", "The port for the API to listen on").Envar("PORT").Default("8080").String()
	APIListenHost := app.Flag("apihost", "The host for the API to listen on").Envar("APILISTENHOST").Default("localhost").String()
	APIAuthUsername := app.Flag("apiuser", "The auth username for the API").Envar("APIAUTHUSERNAME").Required().String()
	APIAuthPass := app.Flag("apipass", "The auth password for the API").Envar("APIAUTHPASSWORD").Required().String()
	LogLevel := app.Flag("loglevel", "The logging level").Envar("LOGLEVEL").Default("WARNING").String()
	ConfigFile := app.Flag("configfile", "The config file").Envar("CONFIGFILE").Default("config.yml").String()
	TestMode := app.Flag("testmode", "Enable for unit testing").Envar("TESTMODE").Default("false").Bool()
	TestURL := app.Flag("testurl", "URL for testing purposes").Envar("TESTURL").Default("").String()
	Bases := app.Flag("bases", "URL for testing purposes").Envar("BASES").StringMap()

	kingpin.MustParse(app.Parse(args))

	o.APIListenPort = *APIListenPort
	o.APIListenHost = *APIListenHost
	o.APIAuthUsername = *APIAuthUsername
	o.APIAuthPass = *APIAuthPass
	o.LogLevel = *LogLevel
	o.ConfigFile = *ConfigFile
	o.TestMode = *TestMode
	o.TestURL = *TestURL
	o.Bases = *Bases

	o.App = app
}
