package apiv1

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/trainer/cli"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/handler"
	"github.com/stretchr/testify/assert"
)

func TestV1_AddHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("AddHandlers registers routes", func(t *testing.T) {
		v := &V1{}
		engine := gin.New()
		o := cli.Options{
			APIAuthUsername: "testuser",
			APIAuthPass:     "testpass",
			LogLevel:        "INFO",
		}
		c := &config.Config{}
		h := &handler.Handler{}

		v.AddHandlers(o, c, engine, h)

		routes := engine.Routes()
		assert.NotEmpty(t, routes)

		foundLaunch := false
		foundRemove := false
		foundStatus := false
		foundConfig := false

		for _, route := range routes {
			if route.Path == "/capi/v1/launch/:plan" && route.Method == "POST" {
				foundLaunch = true
			}
			if route.Path == "/capi/v1/remove" && route.Method == "POST" {
				foundRemove = true
			}
			if route.Path == "/capi/v1/status" && route.Method == "POST" {
				foundStatus = true
			}
			if route.Path == "/capi/v1/config" && route.Method == "POST" {
				foundConfig = true
			}
		}

		assert.True(t, foundLaunch, "Launch route should be registered")
		assert.True(t, foundRemove, "Remove route should be registered")
		assert.True(t, foundStatus, "Status route should be registered")
		assert.True(t, foundConfig, "Config route should be registered")
	})
}

func TestV1_Remove(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Remove endpoint calls handler reset", func(t *testing.T) {
		v := &V1{}
		c := &config.Config{}
		h := &handler.Handler{}
		h.Start()
		defer h.Stop()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		removeHandler := v.Remove(c, h)
		removeHandler(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "remove succeeded", response["message"])
		assert.Equal(t, false, response["error"])
	})
}

func TestV1_Status_NoChange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Status returns no change when no plan is active", func(t *testing.T) {
		v := &V1{}
		c := &config.Config{}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		statusHandler := v.Status(c)
		statusHandler(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "no change", response["message"])
		assert.Equal(t, false, response["error"])
	})
}

func TestV1_ConfigAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ConfigAPI returns YAML configuration", func(t *testing.T) {
		v := &V1{}
		c := &config.Config{
			DefaultPlan: "test_plan",
		}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		configHandler := v.ConfigAPI(c)
		configHandler(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		
		body := w.Body.String()
		assert.NotEmpty(t, body)
		assert.Contains(t, body, "default_plan")
	})
}

func TestV1_Launch_InvalidPlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Launch with invalid plan returns error", func(t *testing.T) {
		v := &V1{}
		c := &config.Config{}

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{
			{Key: "plan", Value: "nonexistent_plan"},
		}

		launchHandler := v.Launch(c)
		launchHandler(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response["error"].(bool))
		assert.NotEmpty(t, response["message"])
	})
}

func TestV1_Authentication(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Routes require authentication", func(t *testing.T) {
		v := &V1{}
		engine := gin.New()
		o := cli.Options{
			APIAuthUsername: "testuser",
			APIAuthPass:     "testpass",
			LogLevel:        "INFO",
		}
		c := &config.Config{}
		h := &handler.Handler{}

		v.AddHandlers(o, c, engine, h)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/capi/v1/status", nil)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Routes accept valid authentication", func(t *testing.T) {
		v := &V1{}
		engine := gin.New()
		o := cli.Options{
			APIAuthUsername: "testuser",
			APIAuthPass:     "testpass",
			LogLevel:        "INFO",
		}
		c := &config.Config{}
		h := &handler.Handler{}

		v.AddHandlers(o, c, engine, h)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/capi/v1/status", nil)
		req.SetBasicAuth("testuser", "testpass")
		engine.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusUnauthorized, w.Code)
	})
}
