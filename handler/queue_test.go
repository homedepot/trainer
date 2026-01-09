package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/trainer/actions"
	"github.com/stretchr/testify/assert"
)

func TestQueue_Add(t *testing.T) {
	gin.SetMode(gin.TestMode)
	qc := &actions.QueueContext{}
	type fields struct {
		rc chan *actions.QueueContext
	}
	type args struct {
		g *actions.QueueContext
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "add",
			fields: fields{},
			args: args{
				g: qc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q.Add(qc)
			select {
			case g := <-q.rc:
				if g != qc {
					t.Errorf("didnt add correct struct somehow")
				}
			default:
				t.Errorf("nothing on channel")
			}
		})
	}
}

func TestQueue_GetUrl_Empty(t *testing.T) {
	qu := &Queue{
		rc: make(chan *actions.QueueContext, 256),
	}

	// Get from empty queue
	retrieved := qu.GetUrl()
	assert.Nil(t, retrieved, "Should return nil for empty queue")
}

func TestQueue_GetUrl_NonBlocking(t *testing.T) {
	qu := &Queue{
		rc: make(chan *actions.QueueContext, 256),
	}

	// Multiple gets should not block
	for i := 0; i < 10; i++ {
		retrieved := qu.GetUrl()
		assert.Nil(t, retrieved, "Should return nil immediately")
	}
}

func TestQueue_MultipleItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	qu := &Queue{
		rc: make(chan *actions.QueueContext, 256),
	}

	// Add multiple items
	items := make([]*actions.QueueContext, 5)
	for i := 0; i < 5; i++ {
		c, _ := gin.CreateTestContext(nil)
		items[i] = &actions.QueueContext{
			Ctx:      c,
			Finished: make(chan bool),
		}
		qu.Add(items[i])
	}

	// Retrieve them in order
	for i := 0; i < 5; i++ {
		retrieved := qu.GetUrl()
		assert.NotNil(t, retrieved, "Should retrieve item %d", i)
		assert.Equal(t, items[i], retrieved, "Should retrieve items in FIFO order")
	}

	// Queue should be empty now
	retrieved := qu.GetUrl()
	assert.Nil(t, retrieved, "Queue should be empty")
}
