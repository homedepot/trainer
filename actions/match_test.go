package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"github.com/homedepot/trainer/structs/plan"
	"github.com/homedepot/trainer/structs/state"
	"github.com/homedepot/trainer/structs/transaction"
	"github.com/mohae/deepcopy"
	"reflect"
	"testing"
)

func TestMatch_Abort(t *testing.T) {
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
			m := &Match{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			y := deepcopy.Copy(m)
			m.Abort()
			if !reflect.DeepEqual(m, y) {
				t.Errorf("Compare() got = %v, want %v", y, m)
				return
			}
		})
	}
}

func TestMatch_GetName(t *testing.T) {
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
			want: "match",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Match{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			if got := m.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatch_LoadFile(t *testing.T) {
	type fields struct {
		Action Action
		Args   ArgStruct
	}
	type args struct {
		n string
		t string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "load json file",
			fields: fields{},
			args: args{
				n: "success.json",
				t: "json",
			},
			want: map[string]interface{}{
				"result": map[string]interface{}{
					"allyourbase.example.com": map[string]interface{}{
						"ci":           "allyourbase.example.com",
						"finished":     float64(1),
						"finishedbool": true,
						"success":      float64(1),
						"successbool":  true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Match{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			got, err := m.LoadFile(tt.args.n, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadFile() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestMatch_LoadString(t *testing.T) {
	type fields struct {
		Action Action
		Args   ArgStruct
	}
	type args struct {
		b string
		t string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "load json file",
			fields: fields{},
			args: args{
				b: success,
				t: "json",
			},
			want: map[string]interface{}{
				"result": map[string]interface{}{
					"allyourbase.example.com": map[string]interface{}{
						"ci":           "allyourbase.example.com",
						"finished":     float64(1),
						"finishedbool": true,
						"success":      float64(1),
						"successbool":  true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Match{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			got, err := m.LoadString(tt.args.b, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadString() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestMatch_Match(t *testing.T) {
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
			name: "successful match",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"match_file":      "success.json",
						"match_file_type": "json",
						"response_type":   "json",
						"variable":        "response",
						"advance_true":    "success",
						"advance_false":   "failure",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"response": success,
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
			name: "successful match, truncated response",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"match_file":      "incomplete.json",
						"match_file_type": "json",
						"response_type":   "json",
						"variable":        "response",
						"advance_true":    "success",
						"advance_false":   "failure",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"response": successtrunc,
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
			name: "unsuccessful match",
			fields: fields{
				Args: ArgStruct{
					Args: map[string]interface{}{
						"match_file":      "success.json",
						"match_file_type": "json",
						"response_type":   "json",
						"variable":        "response",
						"advance_true":    "success",
						"advance_false":   "failure",
					},
				},
			},
			args: args{
				&plan.Plan{
					State: &state.State{
						Variables: map[string]interface{}{
							"response": failure,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Match{
				Action: tt.fields.Action,
				Args:   tt.fields.Args,
			}
			if gotR := m.Execute(tt.args.p); !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("Match() = %v, want %v", gotR, tt.wantR)
			}

		})
	}
}

func TestMatch_SetArgs(t *testing.T) {
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
			m := &Match{
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
