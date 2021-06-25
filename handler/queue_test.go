package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/homedepot/trainer/actions"
	"testing"
)

func TestQueue_Add(t *testing.T) {
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
