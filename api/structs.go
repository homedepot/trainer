package api

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/juju/loggo"
)

// HTTPReturnStruct used to capture
// HTTP request data to return to client.
type HTTPReturnStruct struct {
	Message    string `json:"message"`
	Error      bool   `json:"error"`
	ReturnCode int    `json:"-"`
}

type Version interface {
	AddHandlers(*gin.Engine)
}

// WriteOutput writes HTTP request results to client.
func (e HTTPReturnStruct) WriteOutput(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(e.ReturnCode)

	// Convert the receiver HTTPReturnStruct to JSON.
	outbuf, err := json.Marshal(e)
	if err != nil {
		// Log the marshal error and send a generic error response
		logger := loggo.GetLogger("default")
		logger.Errorf("Failed to marshal HTTPReturnStruct: %v, message: %s", err, e.Message)
		
		// Send a fallback error response
		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte(`{"message":"Internal server error: failed to encode response","error":true}`))
		return
	}
	c.Writer.Write(outbuf)
	logger := loggo.GetLogger("default")
	logger.Warningf("sent API client: %+v", e)
	return
}
