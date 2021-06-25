package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/homedepot/trainer/structs/plan"
	"github.com/homedepot/trainer/structs/state"
	"reflect"
	"testing"
	"time"
)

func TestWait_Abort(t *testing.T) {
	type fields struct {
		Action       Action
		Args         ArgStruct
		Reset        bool
		CurrentTimer int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "abort",
			fields: fields{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wait{
				Action:       tt.fields.Action,
				Args:         tt.fields.Args,
				Reset:        tt.fields.Reset,
				CurrentTimer: tt.fields.CurrentTimer,
			}
			w.Abort()
			if w.Reset != tt.want {
				t.Errorf("Abort() = %v, want %v", w.Reset, tt.want)
			}
		})
	}
}

func TestWait_Execute(t *testing.T) {
	type fields struct {
		Action       Action
		Args         ArgStruct
		Reset        bool
		CurrentTimer int
	}
	type args struct {
		s *plan.Plan
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantR  ExecuteResult
	}{
		{
			name: "execute start",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"duration": 2,
					},
				},
			},
			args: args{
				s: &plan.Plan{
					State: &state.State{},
				},
			},
			wantR: ExecuteResult{
				Err:      nil,
				Complete: false,
				Success:  false,
			},
		}, {
			name: "execute running",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"duration": 2,
					},
				},
			},
			args: args{
				s: &plan.Plan{
					State: &state.State{
						WaitActionStartTime: time.Now().Add(-1 * time.Second),
					},
				},
			},
			wantR: ExecuteResult{
				Err:      nil,
				Complete: false,
				Success:  false,
			},
		}, {
			name: "execute finished",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"duration": 2,
					},
				},
			},
			args: args{
				s: &plan.Plan{
					State: &state.State{
						WaitActionStartTime: time.Now().Add(-2 * time.Second),
					},
				},
			},
			wantR: ExecuteResult{
				Err:      nil,
				Complete: true,
				Success:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wait{
				Action:       tt.fields.Action,
				Args:         tt.fields.Args,
				Reset:        tt.fields.Reset,
				CurrentTimer: tt.fields.CurrentTimer,
			}
			if gotR := w.Execute(tt.args.s); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Execute() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestWait_GetName(t *testing.T) {
	type fields struct {
		Action       Action
		Args         ArgStruct
		Reset        bool
		CurrentTimer int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "getname",
			fields: fields{},
			want:   "wait",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Wait{
				Action:       tt.fields.Action,
				Args:         tt.fields.Args,
				Reset:        tt.fields.Reset,
				CurrentTimer: tt.fields.CurrentTimer,
			}
			if got := w.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
