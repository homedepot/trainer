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
// WriteOutput writes HTTP request results to client.
func (e HTTPReturnStruct) WriteOutput(c *gin.Context) {
	var outbuf = []byte{}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(e.ReturnCode)

	// Convert the receiver HTTPReturnStruct to JSON.
	outbuf, err := json.Marshal(e)
	if err != nil {
		panic("json marshal of HTTPReturnStruct failed ---> fatal application error. " + err.Error() + " " + e.Message)
	}
	c.Writer.Write(outbuf)
	logger := loggo.GetLogger("default")
	logger.Warningf("sent API client: %+v", e)
	return
}
