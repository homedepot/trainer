package router

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/homedepot/trainer/cli"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/handler"
	"github.com/juju/loggo"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		forwardedProto     string
		expectedStatusCode int
		shouldRedirect     bool
		requestPath        string
	}{
		{
			name:               "HTTP request should redirect to HTTPS",
			forwardedProto:     "http",
			expectedStatusCode: http.StatusMovedPermanently,
			shouldRedirect:     true,
			requestPath:        "/test/path",
		},
		{
			name:               "HTTPS request should not redirect",
			forwardedProto:     "https",
			expectedStatusCode: http.StatusOK,
			shouldRedirect:     false,
			requestPath:        "/test/path",
		},
		{
			name:               "No X-Forwarded-Proto header should not redirect",
			forwardedProto:     "",
			expectedStatusCode: http.StatusOK,
			shouldRedirect:     false,
			requestPath:        "/test/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, engine := gin.CreateTestContext(w)

			engine.Use(Redirect())
			engine.GET("/test/path", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", tt.requestPath, nil)
			if tt.forwardedProto != "" {
				req.Header.Set("X-Forwarded-Proto", tt.forwardedProto)
			}
			c.Request = req

			engine.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.shouldRedirect {
				location := w.Header().Get("Location")
				assert.Contains(t, location, "https://")
				assert.NotEmpty(t, location)
			}
		})
	}
}

func TestRouter_CreateRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CreateRouter initializes with nil engine", func(t *testing.T) {
		r := &Router{}
		o := cli.Options{
			APIAuthUsername: "testuser",
			APIAuthPass:     "testpass",
			LogLevel:        "INFO",
		}
		c := &config.Config{}
		logger := loggo.GetLogger("test")
		h := &handler.Handler{}

		r.CreateRouter(o, c, logger, nil, h)

		assert.NotNil(t, r.Router)
	})

	t.Run("CreateRouter uses provided engine", func(t *testing.T) {
		r := &Router{}
		o := cli.Options{
			APIAuthUsername: "testuser",
			APIAuthPass:     "testpass",
			LogLevel:        "INFO",
		}
		c := &config.Config{}
		logger := loggo.GetLogger("test")
		h := &handler.Handler{}
		providedEngine := gin.New()

		r.CreateRouter(o, c, logger, providedEngine, h)

		assert.NotNil(t, r.Router)
		assert.Equal(t, providedEngine, r.Router)
	})

	t.Run("CreateRouter registers routes", func(t *testing.T) {
		r := &Router{}
		o := cli.Options{
			APIAuthUsername: "testuser",
			APIAuthPass:     "testpass",
			LogLevel:        "INFO",
		}
		c := &config.Config{}
		logger := loggo.GetLogger("test")
		h := &handler.Handler{}

		r.CreateRouter(o, c, logger, nil, h)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health-check", nil)
		r.Router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNotFound, w.Code)
	})
}

func TestRouter_AddHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("AddHandler returns a valid handler function", func(t *testing.T) {
		r := &Router{}
		h := &handler.Handler{}

		handlerFunc := r.AddHandler(h)

		assert.NotNil(t, handlerFunc)
	})
}
