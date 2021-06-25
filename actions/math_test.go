package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"errors"
	"fmt"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/homedepot/trainer/structs/state"
	"github.com/mohae/deepcopy"
	"reflect"
	"testing"
)

func TestMath_Abort(t *testing.T) {
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
			w := &Math{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			y := deepcopy.Copy(w)
			w.Abort()
			if !reflect.DeepEqual(w, y) {
				t.Errorf("Compare() got = %v, want %v", y, w)
				return
			}
		})
	}
}

func TestMath_Execute(t *testing.T) {
	type fields struct {
		Action Action
		Args   ArgStruct
	}
	type args struct {
		p *plan.Plan
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantR      ExecuteResult
		wantResult float64
	}{
		{
			name: "adding",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"variable": "source",
						"value":    5,
						"action":   "+",
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
				},
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  true,
				Err:      nil,
			},
			wantResult: 10,
		},
		{
			name: "invalid operation",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"variable": "source",
						"value":    5,
						"action":   "gerfloobity",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"source": float64(5),
						},
					},
				},
			},
			wantR: ExecuteResult{
				Complete: true,
				Success:  false,
				Err:      fmt.Errorf("couldnt execute math operation gerfloobity: %w", errors.New("compare: invalid operation gerfloobity")),
			},
			wantResult: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Math{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			if gotR := w.Execute(tt.args.p); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Execute() = %+v, want %+v", gotR, tt.wantR)
			}
			v, err := tt.args.p.State.GetVariable("source")
			if err != nil {
				t.Errorf("Execute(): error getting variable: %s", err)
			}
			if v != tt.wantResult {
				t.Errorf("Execute() = %v, want %v (%s %s)", v, tt.wantResult, reflect.TypeOf(v).String(), reflect.TypeOf(tt.wantResult).String())
			}
		})
	}
}

func TestMath_GetName(t *testing.T) {
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
			want: "math",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Math{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			if got := w.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMath_SetArgs(t *testing.T) {
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
			m := &Math{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			m.SetArgs(tt.args.i)
			if !reflect.DeepEqual(m.Args.Args, tt.args.i) {
				t.Errorf("SetArgs() = %v, want %v", m.Args.Args, tt.args.i)
			}
		})
	}
}
