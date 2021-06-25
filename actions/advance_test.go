package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"errors"
	"github.com/homedepot/trainer/structs/plan"
	"reflect"
	"testing"
)

func TestAdvance_Execute(t *testing.T) {
	pl, err := Config.FindPlan("basic_test")
	if err != nil {
		t.Errorf("couldn't load config")
		return
	}
	type fields struct {
		Action Action
		Args   ArgStruct
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
			name: "test advance",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"txn": "empty",
					},
				},
			},
			args: args{
				s: pl,
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  true,
				Advance:  true,
				NewTxn:   "empty",
				Err:      nil,
			},
		},
		{
			name: "test advance to invalid transaction",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"txn": "invalid",
					},
				},
			},
			args: args{
				s: pl,
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  false,
				Advance:  false,
				NewTxn:   "",
				Err:      errors.New("no such transaction"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Advance{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			if gotR := a.Execute(tt.args.s); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Execute() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestAdvance_GetName(t *testing.T) {
	type fields struct {
		Action Action
		Args   ArgStruct
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "getname",
			fields: fields{},
			want:   "advance",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Advance{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			if got := a.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
