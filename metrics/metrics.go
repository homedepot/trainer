package metrics

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import "github.com/prometheus/client_golang/prometheus"

/*
Package metrics provides key metrics for application:

Provides metrics for:
- number of changes initiated
- number of changes aborted

The Home Depot
*/

var (
	NumInitiates = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "num_initiates",
		Help: "Number of initiate requests sent",
	})
	NumAborts = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "num_aborts",
		Help: "Number of changes aborted",
	})
)

func Register() {
	prometheus.MustRegister(NumInitiates)
	prometheus.MustRegister(NumAborts)
}
