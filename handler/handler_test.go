package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/homedepot/trainer/structs/state"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Lifecycle(t *testing.T) {
	h := Handler{}
	h.Start()
	
	// Verify channels are initialized
	assert.NotNil(t, h.kill, "kill channel should be initialized")
	assert.NotNil(t, h.abort, "abort channel should be initialized")
	
	time.Sleep(200 * time.Millisecond) // let the ticker tick once
	h.Stop()
}

func TestHandler_Add_PiePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := Handler{}
	// Don't start to avoid ProcessTests

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		URL: &url.URL{
			Path: "/314159265358979323",
		},
	}

	h.Add(c)

	assert.Equal(t, 200, w.Code, "Should return 200 for pi path")
	assert.Equal(t, "I LIKE PIE", w.Body.String(), "Should return pi message")
}

func TestHandler_Add_NoTestInProgress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := Handler{}
	// Don't start handler to avoid ProcessTests being called

	// Ensure no test is in progress
	tst.tst = nil

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		URL: &url.URL{
			Path: "/some/path",
		},
	}

	h.Add(c)

	assert.Equal(t, 500, w.Code, "Should return 500 when no test in progress")
	assert.Equal(t, "no test in progress", w.Body.String())
}

func TestHandler_Add_WithTest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Save original state
	originalTst := tst.tst
	defer func() {
		tst.tst = originalTst
	}()

	// Set up a minimal test with state to allow Add to queue
	tst.tst = &plan.Plan{
		Name: "test",
		State: &state.State{},
	}

	h := Handler{}
	// Note: Not starting handler to keep test simple
	// In real usage, the handler would be started

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		URL: &url.URL{
			Path: "/test/path",
		},
	}

	// Run Add in a goroutine since it will block waiting for Finished signal
	done := make(chan bool)
	go func() {
		h.Add(c)
		done <- true
	}()

	// Small delay to let Add queue the request
	time.Sleep(50 * time.Millisecond)

	// Get the queued context
	qc := q.GetUrl()
	assert.NotNil(t, qc, "Should have queued the request")

	// Signal that we're finished
	qc.Finished <- true

	// Wait for Add to complete
	<-done
}

func TestHandler_Reset(t *testing.T) {
	h := Handler{}
	h.Start()
	defer h.Stop()

	// Set up a test with proper state to avoid nil pointer
	tst.tst = &plan.Plan{
		State: &state.State{},
	}
	tst.processing = false

	// Reset should clear the test
	h.Reset()

	assert.Nil(t, tst.tst, "Test should be nil after reset")
}

func TestHandler_Reset_WaitsForProcessing(t *testing.T) {
	h := Handler{}
	h.Start()
	defer h.Stop()

	// Set up a test that's processing with proper state
	tst.tst = &plan.Plan{
		State: &state.State{},
	}
	tst.processing = true

	// Start reset in a goroutine
	done := make(chan bool)
	go func() {
		h.Reset()
		done <- true
	}()

	// Simulate processing finishing
	time.Sleep(100 * time.Millisecond)
	tst.processing = false

	// Wait for reset to complete
	select {
	case <-done:
		assert.Nil(t, tst.tst, "Test should be nil after reset")
	case <-time.After(5 * time.Second):
		t.Fatal("Reset did not complete in time")
	}
}

func TestLaunchTest(t *testing.T) {
	tests := []struct {
		name        string
		planPath    string
		setupConfig func() *config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name:     "test already in progress",
			planPath: "basic_test",
			setupConfig: func() *config.Config {
				// Set up a test already in progress
				tst.tst = &plan.Plan{Name: "existing"}
				cfg := &config.Config{
					Plans: []plan.Plan{
						{
							Name: "basic_test",
						},
					},
				}
				return cfg
			},
			expectError: true,
			errorMsg:    "test already in progress",
		},
		{
			name:     "plan not found",
			planPath: "nonexistent",
			setupConfig: func() *config.Config {
				tst.tst = nil
				cfg := &config.Config{
					Plans: []plan.Plan{},
				}
				return cfg
			},
			expectError: true,
			errorMsg:    "failed to locate plan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setupConfig()
			err := LaunchTest(cfg, tt.planPath)

			if tt.expectError {
				assert.Error(t, err, "Expected an error")
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "Error message should match")
				}
			} else {
				assert.NoError(t, err, "Should not return an error")
				assert.NotNil(t, tst.tst, "Test should be initialized")
			}

			// Clean up
			tst.tst = nil
		})
	}
}

func TestLaunchTest_Integration(t *testing.T) {
	// This test uses real config file loading
	// Save original state
	originalTst := tst.tst
	defer func() {
		tst.tst = originalTst
	}()

	tst.tst = nil

	// Try to load a real config
	cfg, err := config.NewConfig("../data/config.yml", false, "http://example.com", map[string]string{})
	if err != nil {
		t.Skipf("Skipping integration test - could not load config: %v", err)
		return
	}

	// Find a plan to launch
	if len(cfg.Plans) == 0 {
		t.Skip("No plans available in test config")
		return
	}

	planName := cfg.Plans[0].Name

	err = LaunchTest(cfg, planName)
	if err != nil {
		// If reset failed due to missing transactions, that's expected for some test plans
		if !assert.Contains(t, err.Error(), "transactions", "Expected transaction-related error") {
			t.Logf("LaunchTest failed with: %v", err)
		}
	} else {
		assert.NotNil(t, tst.tst, "Test should be initialized")
		assert.Equal(t, planName, tst.tst.Name, "Plan name should match")
	}

	// Clean up
	tst.tst = nil
}

func TestRemoveTest(t *testing.T) {
	tests := []struct {
		name        string
		setupTest   func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful remove",
			setupTest: func() {
				tst.tst = &plan.Plan{Name: "test"}
			},
			expectError: false,
		},
		{
			name: "no test to remove",
			setupTest: func() {
				tst.tst = nil
			},
			expectError: true,
			errorMsg:    "no test to remove",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupTest()
			err := RemoveTest([16]byte{})

			if tt.expectError {
				assert.Error(t, err, "Expected an error")
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "Error message should match")
				}
			} else {
				assert.NoError(t, err, "Should not return an error")
			}

			// Clean up
			tst.tst = nil
		})
	}
}
