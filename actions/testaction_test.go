package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/homedepot/trainer/structs/plan"
	"reflect"
	"testing"
)

func TestTest_Execute(t1 *testing.T) {
	type fields struct {
		Action Action
		Args   ArgStruct
		Reset  bool
	}
	type args struct {
		p *plan.Plan
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantR  ExecuteResult
	}{
		{
			name: "execute test",
			fields: fields{
				Args: ArgStruct{},
			},
			args: args{
				p: &plan.Plan{},
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  true,
				Err:      nil,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Test{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
				Reset:  tt.fields.Reset,
			}
			if gotR := t.Execute(tt.args.p); !reflect.DeepEqual(gotR, tt.wantR) {
				t1.Errorf("Execute() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
