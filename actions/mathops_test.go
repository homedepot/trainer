package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import "testing"

// There are many, many math operations that I did not test.
// this is because I would basically be testing the go math library,
// and I don't feel like going to that effort for such little gain.

func TestMathOps_Execute(t *testing.T) {
	type fields struct {
		LeftOp     interface{}
		RightOp    interface{}
		leftFloat  float64
		rightFloat float64
		operations map[string]func() float64
	}
	type args struct {
		op string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "add",
			fields: fields{
				LeftOp:  2,
				RightOp: 3,
			},
			args: args{
				"+",
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "incorrect operation",
			fields: fields{
				LeftOp:  2,
				RightOp: 3,
			},
			args: args{
				"(",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MathOps{
				LeftOp:     tt.fields.LeftOp,
				RightOp:    tt.fields.RightOp,
				leftFloat:  tt.fields.leftFloat,
				rightFloat: tt.fields.rightFloat,
				operations: tt.fields.operations,
			}
			got, err := c.Execute(tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
