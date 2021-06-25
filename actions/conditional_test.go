package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"errors"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/homedepot/trainer/structs/state"
	"github.com/homedepot/trainer/structs/transaction"
	"github.com/mohae/deepcopy"
	"reflect"
	"testing"
)

func TestConditional_Abort(t *testing.T) {
	type fields struct {
		Action Action
		Args   ArgStruct
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "everything stays the same",
			fields: fields{
				Args: ArgStruct{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conditional{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			d := deepcopy.Copy(c)
			c.Abort()
			if !reflect.DeepEqual(c, d) {
				t.Errorf("Compare() got = %v, want %v", d, c)
				return
			}
		})
	}
}

func TestConditional_Execute(t *testing.T) {
	type fields struct {
		Action Action
		Args   ArgStruct
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
			name: "comparing equality with value",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"term": map[string]interface{}{
							"variable":          "source",
							"conditional":       "eq",
							"conditional_value": 5,
						},
						"advance_true":  "success",
						"advance_false": "failure",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"source": 5,
						},
					},
					Txn: []transaction.Transaction{
						{
							Name: "success",
						},
						{
							Name: "failure",
						},
					},
				},
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  true,
				Advance:  true,
				NewTxn:   "success",
				Err:      nil,
			},
		},
		{
			name: "comparing inequality with value",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"term": map[string]interface{}{
							"variable":          "source",
							"conditional":       "eq",
							"conditional_value": 5,
						},
						"advance_true":  "success",
						"advance_false": "failure",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"source": 4,
						},
					},
					Txn: []transaction.Transaction{
						{
							Name: "success",
						},
						{
							Name: "failure",
						},
					},
				},
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  false,
				Advance:  true,
				NewTxn:   "failure",
				Err:      nil,
			},
		},
		{
			name: "comparing equality with variable",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"term": map[string]interface{}{
							"variable":        "source",
							"conditional":     "eq",
							"conditional_var": "dest",
						},
						"advance_true":  "success",
						"advance_false": "failure",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"source": 5,
							"dest":   5,
						},
					},
					Txn: []transaction.Transaction{
						{
							Name: "success",
						},
						{
							Name: "failure",
						},
					},
				},
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  true,
				Advance:  true,
				NewTxn:   "success",
				Err:      nil,
			},
		},
		{
			name: "comparing inequality with variable",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"term": map[string]interface{}{
							"variable":        "source",
							"conditional":     "eq",
							"conditional_var": "dest",
						},
						"advance_true":  "success",
						"advance_false": "failure",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"source": 4,
							"dest":   5,
						},
					},
					Txn: []transaction.Transaction{
						{
							Name: "success",
						},
						{
							Name: "failure",
						},
					},
				},
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  false,
				Advance:  true,
				NewTxn:   "failure",
				Err:      nil,
			},
		},
		{
			name: "comparing with invalid operation",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"term": map[string]interface{}{
							"variable":        "source",
							"conditional":     "haha",
							"conditional_var": "dest",
						},
						"advance_true":  "success",
						"advance_false": "failure",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"source": 4,
							"dest":   5,
						},
					},
					Txn: []transaction.Transaction{
						{
							Name: "success",
						},
						{
							Name: "failure",
						},
					},
				},
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  false,
				Advance:  false,
				NewTxn:   "",
				Err:      errors.New("compare: invalid operator haha"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conditional{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			if gotR := c.Execute(tt.args.p); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Execute() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestConditional_GetName(t *testing.T) {
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
			name: "get name",
			fields: fields{
				Args: ArgStruct{},
			},
			want: "conditional",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conditional{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			if got := c.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditional_SetArgs(t *testing.T) {
	a := map[string]interface{}{
		"arg1": "anarg",
	}

	type fields struct {
		Action Action
		Args   ArgStruct
	}
	type args struct {
		i map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set an argument",
			fields: fields{
				Args: ArgStruct{},
			},
			args: args{
				i: a,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conditional{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			c.SetArgs(tt.args.i)
			if !reflect.DeepEqual(c.Args.Args, tt.args.i) {
				t.Errorf("SetArgs() = %v, want %v", c.Args.Args, tt.args.i)
			}
		})
	}
}
