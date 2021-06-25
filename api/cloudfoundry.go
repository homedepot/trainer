package api

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/loggo"
)

// CloudFoundry facilitates CloudFoundry calls, which occur at
// varied times. If we do not catch these, it will appear as if
// something is in the application, and we need to avoid this.
func CloudFoundry(c *gin.Context) {
	logger := loggo.GetLogger("default")

	logger.Tracef("Options request from PCF ignored.")
}
