package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/homedepot/trainer/actions"
)

var q *Queue

type Queue struct {
	// contains a list of queued entries.
	rc chan *actions.QueueContext
}

func init() {
	q = &Queue{}
	q.rc = make(chan *actions.QueueContext, 256)
}

func (qu *Queue) Add(g *actions.QueueContext) {
	qu.rc <- g
}

func (qu *Queue) GetUrl() *actions.QueueContext {
	select {
	case g := <-qu.rc:
		return g
	default:
		return nil
	}
}
