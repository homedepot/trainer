package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/gin-gonic/gin"
	"github.com/homedepot/trainer/actions"
	"github.com/juju/loggo"
	"time"
)

type Handler struct {
	kill  chan *bool
	abort chan *bool
}

func (h *Handler) Start() {
	logger := loggo.GetLogger("default")
	c := make(chan *bool)
	abort := make(chan *bool, 1)
	go func(c chan *bool) {
		logger.Warningf("Starting runner...")
		Runner(c, abort)
	}(c)
	h.kill = c
	h.abort = abort
}

func (h *Handler) Add(c *gin.Context) {
	logger := loggo.GetLogger("default")
	if c.Request.URL.Path == "/314159265358979323" {
		c.Writer.WriteHeader(200)
		c.Writer.Write([]byte("I LIKE PIE"))
		return
	}

	if tst.tst == nil {
		c.Writer.WriteHeader(500)
		c.Writer.Write([]byte("no test in progress"))
		return
	}

	qc := &actions.QueueContext{
		Ctx:      c,
		Finished: make(chan bool),
	}
	q.Add(qc)

	logger.Tracef("waiting for request to finish...")
	select {
	case <-qc.Finished:
		return
	}
}

func (h *Handler) Stop() {
	logger := loggo.GetLogger("default")
	logger.Warningf("Stopping Runner")
	kill := true
	h.kill <- &kill
}

func (h *Handler) Reset() {
	logger := loggo.GetLogger("default")
	for {
		if tst.processing == true {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		break
	}
	logger.Tracef("resetting...")
	abort := true
	h.abort <- &abort
	// TODO:  need a better way to make sure the abort has completed, this is terrible.
	logger.Tracef("reset.")
	time.Sleep(2 * time.Second)
	tst.tst = nil
}
