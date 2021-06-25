package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"testing"
	"time"
)

/*func TestHandler_Add(t *testing.T) {
	type fields struct {
		kill chan bool
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
	}{
		{
			name:   "add an entry",
			fields: fields{},
			args: args{
				c: &gin.Context{
					Request: &http.Request{
						URL: &url.URL{
							Path: "/a/url",
						},
					},
				},
			},
			want:    "/a/url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			h.Add(tt.args.c)

			turl := q.GetUrl()
			if turl == nil {
				t.Errorf("Add didn't succeed")
			}
		})
	}
}*/

func TestHandler_Lifecycle(t *testing.T) {
	// This test will always pass, but the error messages at least show whether the runner is starting
	// and stopping,
	h := Handler{}
	h.Start()
	time.Sleep(200 * time.Millisecond) // let the ticker tick once
	h.Stop()
}
