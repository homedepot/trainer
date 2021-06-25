package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import "testing"

func TestConditionalOps_Compare(t *testing.T) {
	type fields struct {
		LeftOp   interface{}
		RightOp  interface{}
		cmpFuncs map[string]func() (bool, error)
	}
	type args struct {
		op string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionalOps{
				LeftOp:   tt.fields.LeftOp,
				RightOp:  tt.fields.RightOp,
				cmpFuncs: tt.fields.cmpFuncs,
			}
			got, err := c.Compare(tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Compare() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalOps_CompareEq(t *testing.T) {
	c := NewConditionalOps()

	type fields struct {
		LeftOp   interface{}
		RightOp  interface{}
		cmpFuncs map[string]func() (bool, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "equal string",
			fields: fields{
				LeftOp:   "astring",
				RightOp:  "astring",
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "unequal string",
			fields: fields{
				LeftOp:   "astring",
				RightOp:  "astringent",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "equal bool true",
			fields: fields{
				LeftOp:   true,
				RightOp:  true,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "equal bool false",
			fields: fields{
				LeftOp:   false,
				RightOp:  false,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "unequal bool 1",
			fields: fields{
				LeftOp:   true,
				RightOp:  false,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "unequal bool 2",
			fields: fields{
				LeftOp:   false,
				RightOp:  true,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "equal float64",
			fields: fields{
				LeftOp:   3.14159,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "unequal float64",
			fields: fields{
				LeftOp:   3.14159,
				RightOp:  2.71828,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "equal int",
			fields: fields{
				LeftOp:   3,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "unequal float64",
			fields: fields{
				LeftOp:   3,
				RightOp:  2,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionalOps{
				LeftOp:   tt.fields.LeftOp,
				RightOp:  tt.fields.RightOp,
				cmpFuncs: tt.fields.cmpFuncs,
			}
			got, err := c.CompareEq()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareEq() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareEq() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalOps_CompareGe(t *testing.T) {
	c := NewConditionalOps()

	type fields struct {
		LeftOp   interface{}
		RightOp  interface{}
		cmpFuncs map[string]func() (bool, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "equal string",
			fields: fields{
				LeftOp:   "as",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "less string",
			fields: fields{
				LeftOp:   "ar",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "greater string",
			fields: fields{
				LeftOp:   "at",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "boolean should fail",
			fields: fields{
				LeftOp:   true,
				RightOp:  true,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "equal float64",
			fields: fields{
				LeftOp:   3.14159,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "less float64",
			fields: fields{
				LeftOp:   3.14159,
				RightOp:  3.14160,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "greater float64",
			fields: fields{
				LeftOp:   3.14160,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "equal int",
			fields: fields{
				LeftOp:   3,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "less int",
			fields: fields{
				LeftOp:   2,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "greater int",
			fields: fields{
				LeftOp:   4,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "float64 and int equal",
			fields: fields{
				LeftOp:   3.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "float64 and int less",
			fields: fields{
				LeftOp:   2.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "float64 and int greater",
			fields: fields{
				LeftOp:   4.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "float64 and string should fail",
			fields: fields{
				LeftOp:   4.00,
				RightOp:  "3",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionalOps{
				LeftOp:   tt.fields.LeftOp,
				RightOp:  tt.fields.RightOp,
				cmpFuncs: tt.fields.cmpFuncs,
			}
			got, err := c.CompareGe()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareGe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareGe() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalOps_CompareGt(t *testing.T) {
	c := NewConditionalOps()

	type fields struct {
		LeftOp   interface{}
		RightOp  interface{}
		cmpFuncs map[string]func() (bool, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "equal string",
			fields: fields{
				LeftOp:   "as",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "less string",
			fields: fields{
				LeftOp:   "ar",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "greater string",
			fields: fields{
				LeftOp:   "at",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "boolean should fail",
			fields: fields{
				LeftOp:   true,
				RightOp:  true,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "equal float64",
			fields: fields{
				LeftOp:   3.14159,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "less float64",
			fields: fields{
				LeftOp:   3.14158,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "greater float64",
			fields: fields{
				LeftOp:   3.14160,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "equal int",
			fields: fields{
				LeftOp:   3,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "less int",
			fields: fields{
				LeftOp:   2,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "greater int",
			fields: fields{
				LeftOp:   4,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "float64 and int equal",
			fields: fields{
				LeftOp:   3.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "float64 and int less",
			fields: fields{
				LeftOp:   2.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "float64 and int greater",
			fields: fields{
				LeftOp:   4.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "float64 and string should fail",
			fields: fields{
				LeftOp:   4.00,
				RightOp:  "3",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionalOps{
				LeftOp:   tt.fields.LeftOp,
				RightOp:  tt.fields.RightOp,
				cmpFuncs: tt.fields.cmpFuncs,
			}
			got, err := c.CompareGt()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareGt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareGt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalOps_CompareLe(t *testing.T) {
	c := NewConditionalOps()

	type fields struct {
		LeftOp   interface{}
		RightOp  interface{}
		cmpFuncs map[string]func() (bool, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "equal string",
			fields: fields{
				LeftOp:   "as",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "less string",
			fields: fields{
				LeftOp:   "ar",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "greater string",
			fields: fields{
				LeftOp:   "at",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "boolean should fail",
			fields: fields{
				LeftOp:   true,
				RightOp:  true,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "equal float64",
			fields: fields{
				LeftOp:   3.1415926,
				RightOp:  3.1415926,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "less float64",
			fields: fields{
				LeftOp:   3.1415925,
				RightOp:  3.1415926,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "greater float64",
			fields: fields{
				LeftOp:   3.1415927,
				RightOp:  3.1415926,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "equal int",
			fields: fields{
				LeftOp:   3,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "less int",
			fields: fields{
				LeftOp:   2,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "greater int",
			fields: fields{
				LeftOp:   4,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "float64 and int equal",
			fields: fields{
				LeftOp:   3.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "float64 and int less",
			fields: fields{
				LeftOp:   2.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "float64 and int greater",
			fields: fields{
				LeftOp:   4.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "float64 and string should fail",
			fields: fields{
				LeftOp:   4.00,
				RightOp:  "3",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionalOps{
				LeftOp:   tt.fields.LeftOp,
				RightOp:  tt.fields.RightOp,
				cmpFuncs: tt.fields.cmpFuncs,
			}
			got, err := c.CompareLe()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareLe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareLe() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalOps_CompareLt(t *testing.T) {
	c := NewConditionalOps()

	type fields struct {
		LeftOp   interface{}
		RightOp  interface{}
		cmpFuncs map[string]func() (bool, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "equal string",
			fields: fields{
				LeftOp:   "as",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "less string",
			fields: fields{
				LeftOp:   "ar",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "greater string",
			fields: fields{
				LeftOp:   "at",
				RightOp:  "as",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "boolean should fail",
			fields: fields{
				LeftOp:   true,
				RightOp:  true,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "equal float64",
			fields: fields{
				LeftOp:   3.14159,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "less float64",
			fields: fields{
				LeftOp:   3.14158,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "greater float64",
			fields: fields{
				LeftOp:   3.14160,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "equal int",
			fields: fields{
				LeftOp:   3,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "less int",
			fields: fields{
				LeftOp:   2,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "greater int",
			fields: fields{
				LeftOp:   4,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "float64 and int equal",
			fields: fields{
				LeftOp:   3.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "float64 and int less",
			fields: fields{
				LeftOp:   2.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "float64 and int greater",
			fields: fields{
				LeftOp:   4.00,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "float64 and string should fail",
			fields: fields{
				LeftOp:   4.00,
				RightOp:  "3",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionalOps{
				LeftOp:   tt.fields.LeftOp,
				RightOp:  tt.fields.RightOp,
				cmpFuncs: tt.fields.cmpFuncs,
			}
			got, err := c.CompareLt()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareLt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareLt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalOps_CompareNe(t *testing.T) {
	c := NewConditionalOps()
	type fields struct {
		LeftOp   interface{}
		RightOp  interface{}
		cmpFuncs map[string]func() (bool, error)
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "equal string",
			fields: fields{
				LeftOp:   "astring",
				RightOp:  "astring",
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "unequal string",
			fields: fields{
				LeftOp:   "astring",
				RightOp:  "astringent",
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "equal bool true",
			fields: fields{
				LeftOp:   true,
				RightOp:  true,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "equal bool false",
			fields: fields{
				LeftOp:   false,
				RightOp:  false,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "unequal bool 1",
			fields: fields{
				LeftOp:   true,
				RightOp:  false,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "unequal bool 2",
			fields: fields{
				LeftOp:   false,
				RightOp:  true,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "equal float64",
			fields: fields{
				LeftOp:   3.14159,
				RightOp:  3.14159,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "unequal float64",
			fields: fields{
				LeftOp:   3.14159,
				RightOp:  2.71828,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "equal int",
			fields: fields{
				LeftOp:   3,
				RightOp:  3,
				cmpFuncs: c.cmpFuncs,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "unequal float64",
			fields: fields{
				LeftOp:   3,
				RightOp:  2,
				cmpFuncs: c.cmpFuncs,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConditionalOps{
				LeftOp:   tt.fields.LeftOp,
				RightOp:  tt.fields.RightOp,
				cmpFuncs: tt.fields.cmpFuncs,
			}
			got, err := c.CompareNe()
			if (err != nil) != tt.wantErr {
				t.Errorf("CompareNe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompareNe() got = %v, want %v", got, tt.want)
			}
		})
	}
}
