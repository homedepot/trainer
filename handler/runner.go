package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/juju/loggo"
	"time"
)

func Runner(b chan *bool, abort chan *bool) {
	logger := loggo.GetLogger("default")
	pticker := time.NewTicker(200 * time.Millisecond)
	pTickChan := pticker.C
	for {
		select {
		case <-b:
			logger.Warningf("exiting Runner")
			pticker.Stop()
			return
		case <-pTickChan:
			ProcessTests(abort)
		}
	}
}
