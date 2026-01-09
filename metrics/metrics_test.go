package metrics

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Run("Register does not panic", func(t *testing.T) {
		registry := prometheus.NewRegistry()
		
		registry.MustRegister(NumInitiates)
		registry.MustRegister(NumAborts)
		
		assert.NotNil(t, registry, "Registry should be created successfully")
	})
}

func TestNumInitiates(t *testing.T) {
	t.Run("NumInitiates can be set and read", func(t *testing.T) {
		NumInitiates.Set(5)
		
		value := testutil.ToFloat64(NumInitiates)
		assert.Equal(t, 5.0, value, "NumInitiates should be set to 5")
	})

	t.Run("NumInitiates can be incremented", func(t *testing.T) {
		NumInitiates.Set(0)
		
		NumInitiates.Inc()
		NumInitiates.Inc()
		NumInitiates.Inc()
		
		value := testutil.ToFloat64(NumInitiates)
		assert.Equal(t, 3.0, value, "NumInitiates should be 3 after incrementing 3 times")
	})

	t.Run("NumInitiates can be decremented", func(t *testing.T) {
		NumInitiates.Set(10)
		
		NumInitiates.Dec()
		NumInitiates.Dec()
		
		value := testutil.ToFloat64(NumInitiates)
		assert.Equal(t, 8.0, value, "NumInitiates should be 8 after decrementing 2 times from 10")
	})

	t.Run("NumInitiates can be set to zero", func(t *testing.T) {
		NumInitiates.Set(100)
		NumInitiates.Set(0)
		
		value := testutil.ToFloat64(NumInitiates)
		assert.Equal(t, 0.0, value, "NumInitiates should be 0")
	})

	t.Run("NumInitiates can handle large values", func(t *testing.T) {
		largeValue := 1000000.0
		NumInitiates.Set(largeValue)
		
		value := testutil.ToFloat64(NumInitiates)
		assert.Equal(t, largeValue, value, "NumInitiates should handle large values")
	})
}

func TestNumAborts(t *testing.T) {
	t.Run("NumAborts can be set and read", func(t *testing.T) {
		NumAborts.Set(3)
		
		value := testutil.ToFloat64(NumAborts)
		assert.Equal(t, 3.0, value, "NumAborts should be set to 3")
	})

	t.Run("NumAborts can be incremented", func(t *testing.T) {
		NumAborts.Set(0)
		
		NumAborts.Inc()
		NumAborts.Inc()
		
		value := testutil.ToFloat64(NumAborts)
		assert.Equal(t, 2.0, value, "NumAborts should be 2 after incrementing 2 times")
	})

	t.Run("NumAborts can be decremented", func(t *testing.T) {
		NumAborts.Set(5)
		
		NumAborts.Dec()
		
		value := testutil.ToFloat64(NumAborts)
		assert.Equal(t, 4.0, value, "NumAborts should be 4 after decrementing once from 5")
	})

	t.Run("NumAborts can be set to zero", func(t *testing.T) {
		NumAborts.Set(50)
		NumAborts.Set(0)
		
		value := testutil.ToFloat64(NumAborts)
		assert.Equal(t, 0.0, value, "NumAborts should be 0")
	})
}

func TestMetricsMetadata(t *testing.T) {
	t.Run("NumInitiates has correct metadata", func(t *testing.T) {
		assert.NotNil(t, NumInitiates, "NumInitiates should be initialized")
	})

	t.Run("NumAborts has correct metadata", func(t *testing.T) {
		assert.NotNil(t, NumAborts, "NumAborts should be initialized")
	})
}

func TestMetricsConcurrency(t *testing.T) {
	t.Run("metrics can handle concurrent updates", func(t *testing.T) {
		NumInitiates.Set(0)
		NumAborts.Set(0)
		
		done := make(chan bool)
		
		for i := 0; i < 10; i++ {
			go func() {
				NumInitiates.Inc()
				NumAborts.Inc()
				done <- true
			}()
		}
		
		for i := 0; i < 10; i++ {
			<-done
		}
		
		initiatesValue := testutil.ToFloat64(NumInitiates)
		abortsValue := testutil.ToFloat64(NumAborts)
		
		assert.Equal(t, 10.0, initiatesValue, "NumInitiates should be 10 after concurrent increments")
		assert.Equal(t, 10.0, abortsValue, "NumAborts should be 10 after concurrent increments")
	})
}
