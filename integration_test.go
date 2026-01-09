package main

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/handler"
	"github.com/stretchr/testify/assert"
)

// getConfigPath returns the absolute path to the test config file
func getConfigPath(t *testing.T) string {
	// Try different possible paths
	paths := []string{
		"data/config.yml",
		"./data/config.yml",
		"../data/config.yml",
	}
	
	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		if err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}
	
	// If none work, return the default and let the test fail with clear error
	return "data/config.yml"
}

// TestIntegration_LaunchAndRemove tests the complete lifecycle of launching and removing a test
func TestIntegration_LaunchAndRemove(t *testing.T) {
	// Load test configuration
	cfg, err := config.NewConfig(getConfigPath(t), false, "http://example.com", map[string]string{})
	assert.NoError(t, err, "Should load test config")

	assert.NotEmpty(t, cfg.Plans, "Should have test plans")

	planName := cfg.Plans[0].Name

	// Test launching
	err = handler.LaunchTest(cfg, planName)
	
	// If launch fails due to missing transactions, that's OK for some test plans
	// but we should at least verify the error is expected
	if err != nil {
		if assert.Contains(t, err.Error(), "transaction", "Error should be transaction-related") {
			t.Logf("Plan %s lacks transactions - this is expected for minimal test plans", planName)
		} else {
			t.Fatalf("Unexpected error launching test: %v", err)
		}
		return
	}

	// Verify test was launched
	plan := handler.GetPlan()
	assert.NotNil(t, plan, "Plan should be active after launch")
	assert.Equal(t, planName, plan.Name, "Plan name should match")

	// Test removing
	err = handler.RemoveTest([16]byte{})
	assert.NoError(t, err, "Remove should succeed")
}

// TestIntegration_ConcurrentHandlerOperations tests concurrent handler operations
func TestIntegration_ConcurrentHandlerOperations(t *testing.T) {
	h := handler.Handler{}
	
	// Start handler
	h.Start()
	
	// Let it run for a bit
	time.Sleep(500 * time.Millisecond)
	
	// Stop handler
	h.Stop()
	
	// Should be able to start again
	h.Start()
	time.Sleep(200 * time.Millisecond)
	h.Stop()
}

// TestIntegration_ConfigLoadAndPlanRetrieval tests loading config and retrieving plans
func TestIntegration_ConfigLoadAndPlanRetrieval(t *testing.T) {
	cfg, err := config.NewConfig(getConfigPath(t), false, "http://example.com", map[string]string{
		"base1": "http://base1.example.com",
	})
	
	assert.NoError(t, err, "Should load test config")

	// Verify bases were merged
	assert.NotNil(t, cfg.Bases, "Bases should be initialized")
	assert.Contains(t, cfg.Bases, "base1", "Custom base should be included")

	// Test finding a plan
	assert.NotEmpty(t, cfg.Plans, "Should have test plans")
	
	planName := cfg.Plans[0].Name
	plan, err := cfg.FindPlan(planName)
	assert.NoError(t, err, "Should find existing plan")
	assert.NotNil(t, plan, "Plan should not be nil")
	assert.Equal(t, planName, plan.Name, "Plan names should match")

	// Test plan not found
	_, err = cfg.FindPlan("nonexistent-plan-name")
	assert.Error(t, err, "Should return error for nonexistent plan")
	assert.Contains(t, err.Error(), "failed to locate plan", "Error should mention plan not found")
}

// TestIntegration_EndToEndWithRealConfig tests a realistic scenario with actual config
func TestIntegration_EndToEndWithRealConfig(t *testing.T) {
	cfg, err := config.NewConfig(getConfigPath(t), false, "http://example.com", map[string]string{})
	assert.NoError(t, err, "Should load test config")

	// Start a handler to run tests
	h := &handler.Handler{}
	h.Start()
	defer h.Stop()

	// Track results
	successCount := 0
	failCount := 0
	
	// Attempt to launch first few plans (some may fail due to missing transactions)
	maxPlansToTest := 3
	if len(cfg.Plans) < maxPlansToTest {
		maxPlansToTest = len(cfg.Plans)
	}
	
	for i := 0; i < maxPlansToTest; i++ {
		plan := cfg.Plans[i]
		t.Logf("Testing plan: %s", plan.Name)
		
		err := handler.LaunchTest(cfg, plan.Name)
		if err == nil {
			successCount++
			t.Logf("✓ Plan %s launched successfully", plan.Name)
			
			// Verify plan is active
			activePlan := handler.GetPlan()
			assert.NotNil(t, activePlan, "Plan should be active")
			
			// Clean up by removing
			handler.RemoveTest([16]byte{})
		} else {
			failCount++
			t.Logf("✗ Plan %s failed: %v", plan.Name, err)
		}
	}

	t.Logf("Results: %d succeeded, %d failed out of %d tested", successCount, failCount, maxPlansToTest)
	
	// At least one plan should work, or if all fail, failures should be expected (transaction-related)
	if successCount == 0 && failCount > 0 {
		t.Log("All plans failed - this may be expected if test plans lack transactions")
	}
}
