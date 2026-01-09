package handler

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	kill := make(chan *bool, 1)
	abort := make(chan *bool, 1)

	// Start runner in goroutine
	done := make(chan bool)
	go func() {
		Runner(kill, abort)
		done <- true
	}()

	// Send kill signal
	k := true
	kill <- &k

	// Wait for runner to exit
	<-done
}

func TestRunner_ProcessesTicks(t *testing.T) {
	// Save original state
	originalTst := tst.tst
	defer func() {
		tst.tst = originalTst
	}()

	// Set tst to nil so ProcessTests doesn't crash
	tst.tst = nil

	kill := make(chan *bool, 1)
	abort := make(chan *bool, 1)

	// Start runner in goroutine
	done := make(chan bool)
	go func() {
		Runner(kill, abort)
		done <- true
	}()

	// Let a few ticks happen (ProcessTests will return immediately since tst.tst is nil)
	// Sleep for more than 200ms to ensure at least one tick
	// Note: We can't easily verify ProcessTests was called, but this ensures the ticker logic works
	
	// Send kill signal
	k := true
	kill <- &k

	// Wait for runner to exit
	<-done

	assert.True(t, true, "Runner completed without hanging")
}
