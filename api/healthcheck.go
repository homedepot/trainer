package api

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {

	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte("Looks good to me"))
}
