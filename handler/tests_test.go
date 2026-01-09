package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"testing"

	"github.com/homedepot/trainer/structs/plan"
	"github.com/homedepot/trainer/structs/planaction"
	"github.com/homedepot/trainer/structs/state"
	"github.com/homedepot/trainer/structs/transaction"
	"github.com/stretchr/testify/assert"
)

func TestFindSatisfaction(t *testing.T) {
	pa1 := &planaction.PlanAction{Type: "test1"}
	pa2 := &planaction.PlanAction{Type: "test2"}
	
	groups := []*SatisfyGroup{
		{
			name:   "group1",
			action: []*planaction.PlanAction{pa1},
		},
		{
			name:   "group2",
			action: []*planaction.PlanAction{pa2},
		},
	}

	tests := []struct {
		name       string
		searchName string
		groups     []*SatisfyGroup
		expectNil  bool
	}{
		{
			name:       "find existing group",
			searchName: "group1",
			groups:     groups,
			expectNil:  false,
		},
		{
			name:       "group not found",
			searchName: "nonexistent",
			groups:     groups,
			expectNil:  true,
		},
		{
			name:       "empty groups",
			searchName: "test",
			groups:     []*SatisfyGroup{},
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSatisfaction(tt.searchName, tt.groups)
			if tt.expectNil {
				assert.Nil(t, result, "Should return nil")
			} else {
				assert.NotNil(t, result, "Should find the group")
				assert.Equal(t, tt.searchName, result.name, "Should return correct group")
			}
		})
	}
}

func TestCollateActions(t *testing.T) {
	tests := []struct {
		name            string
		actions         []planaction.PlanAction
		expectedGroups  int
		expectedInFirst int
	}{
		{
			name: "no satisfy groups",
			actions: []planaction.PlanAction{
				{Type: "test1", SatisfyGroup: ""},
				{Type: "test2", SatisfyGroup: ""},
			},
			expectedGroups:  2,
			expectedInFirst: 1,
		},
		{
			name: "one satisfy group",
			actions: []planaction.PlanAction{
				{Type: "test1", SatisfyGroup: "group1"},
				{Type: "test2", SatisfyGroup: "group1"},
			},
			expectedGroups:  1,
			expectedInFirst: 2,
		},
		{
			name: "mixed satisfy groups",
			actions: []planaction.PlanAction{
				{Type: "test1", SatisfyGroup: ""},
				{Type: "test2", SatisfyGroup: "group1"},
				{Type: "test3", SatisfyGroup: "group1"},
				{Type: "test4", SatisfyGroup: "group2"},
			},
			expectedGroups:  3,
			expectedInFirst: 1,
		},
		{
			name:            "empty actions",
			actions:         []planaction.PlanAction{},
			expectedGroups:  0,
			expectedInFirst: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CollateActions(tt.actions)
			assert.Equal(t, tt.expectedGroups, len(result), "Should create correct number of groups")
			if tt.expectedGroups > 0 {
				assert.Equal(t, tt.expectedInFirst, len(result[0].action), "First group should have correct number of actions")
			}
		})
	}
}

func TestGetPlan(t *testing.T) {
	// Save original state
	originalTst := tst.tst
	defer func() {
		tst.tst = originalTst
	}()

	tests := []struct {
		name      string
		setupPlan func()
		expectNil bool
	}{
		{
			name: "plan exists",
			setupPlan: func() {
				tst.tst = &plan.Plan{Name: "test"}
			},
			expectNil: false,
		},
		{
			name: "no plan",
			setupPlan: func() {
				tst.tst = nil
			},
			expectNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupPlan()
			result := GetPlan()
			if tt.expectNil {
				assert.Nil(t, result, "Should return nil")
			} else {
				assert.NotNil(t, result, "Should return plan")
			}
		})
	}
}

func TestParseStringTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		plan     *plan.Plan
		expected string
	}{
		{
			name:  "simple variable substitution",
			input: "Hello <<.Variables.name>>",
			plan: &plan.Plan{
				State: &state.State{
					Variables: map[string]interface{}{
						"name": "World",
					},
				},
				Bases: map[string]string{},
			},
			expected: "Hello World",
		},
		{
			name:  "base URL substitution",
			input: "URL: <<.Bases.api>>",
			plan: &plan.Plan{
				State: &state.State{
					Variables: map[string]interface{}{},
				},
				Bases: map[string]string{
					"api": "https://api.example.com",
				},
			},
			expected: "URL: https://api.example.com",
		},
		{
			name:  "non-string variable ignored",
			input: "Count: <<.Variables.count>>",
			plan: &plan.Plan{
				State: &state.State{
					Variables: map[string]interface{}{
						"count": 42,
					},
				},
				Bases: map[string]string{},
			},
			expected: "Count: <no value>",
		},
		{
			name:  "template with Now timestamp",
			input: "Time: <<.Now>>",
			plan: &plan.Plan{
				State: &state.State{
					Variables: map[string]interface{}{},
				},
				Bases: map[string]string{},
			},
			expected: "Time: ", // Will have timestamp appended
		},
		{
			name:  "invalid template returns original",
			input: "Bad template <<.Missing",
			plan: &plan.Plan{
				State: &state.State{
					Variables: map[string]interface{}{},
				},
				Bases: map[string]string{},
			},
			expected: "Bad template <<.Missing",
		},
		{
			name:  "multiple variables",
			input: "<<.Variables.first>> and <<.Variables.second>>",
			plan: &plan.Plan{
				State: &state.State{
					Variables: map[string]interface{}{
						"first":  "foo",
						"second": "bar",
					},
				},
				Bases: map[string]string{},
			},
			expected: "foo and bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseStringTemplate(tt.plan, tt.input)
			if tt.name == "template with Now timestamp" {
				// Just verify it starts with the expected prefix
				assert.Contains(t, result, "Time: ", "Should contain time prefix")
			} else {
				assert.Equal(t, tt.expected, result, "Template parsing should match expected")
			}
		})
	}
}

func TestProcessTests_NilTest(t *testing.T) {
	// Save original state
	originalTst := tst.tst
	originalProcessing := tst.processing
	defer func() {
		tst.tst = originalTst
		tst.processing = originalProcessing
	}()

	tst.tst = nil
	abort := make(chan *bool, 1)

	ProcessTests(abort)

	assert.False(t, tst.processing, "Should not be processing after nil test")
}

func TestProcessTests_StopVar(t *testing.T) {
	// Save original state
	originalTst := tst.tst
	originalProcessing := tst.processing
	defer func() {
		tst.tst = originalTst
		tst.processing = originalProcessing
	}()

	tst.tst = &plan.Plan{
		StopVar: "stop",
		State: &state.State{
			Variables: map[string]interface{}{
				"stop": true,
			},
			States: []state.StateEntry{
				{Status: "running"},
			},
		},
	}
	abort := make(chan *bool, 1)

	ProcessTests(abort)

	assert.Equal(t, "stopped", tst.tst.State.States[0].Status, "Should mark as stopped")
	assert.False(t, tst.processing, "Should not be processing after stop")
}

func TestProcessTests_StateError(t *testing.T) {
	// Save original state
	originalTst := tst.tst
	originalProcessing := tst.processing
	defer func() {
		tst.tst = originalTst
		tst.processing = originalProcessing
	}()

	tst.tst = &plan.Plan{
		State: &state.State{
			Err: assert.AnError,
		},
	}
	abort := make(chan *bool, 1)

	ProcessTests(abort)

	assert.False(t, tst.processing, "Should not be processing when state has error")
}

func TestProcessTests_InvalidTransaction(t *testing.T) {
	// Save original state
	originalTst := tst.tst
	originalProcessing := tst.processing
	defer func() {
		tst.tst = originalTst
		tst.processing = originalProcessing
	}()

	tst.tst = &plan.Plan{
		State: &state.State{
			Transaction: "nonexistent",
			States: []state.StateEntry{
				{Status: "running"},
			},
		},
		Txn: []transaction.Transaction{},
	}
	abort := make(chan *bool, 1)

	ProcessTests(abort)

	assert.NotNil(t, tst.tst.State.Err, "Should set error for invalid transaction")
	assert.Contains(t, tst.tst.State.Err.Error(), "invalid transaction", "Error should mention invalid transaction")
	assert.False(t, tst.processing, "Should not be processing after error")
}
