package core

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlackPost_Success(t *testing.T) {
	t.Run("successful POST returns no error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "no-cache", r.Header.Get("Cache-Control"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ok": true}`))
		}))
		defer server.Close()

		payload := []byte(`{"text": "test message"}`)
		err := SlackPost(payload, server.URL, false)

		assert.NoError(t, err)
	})

	t.Run("201 Created is successful", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		}))
		defer server.Close()

		payload := []byte(`{"text": "test"}`)
		err := SlackPost(payload, server.URL, false)

		assert.NoError(t, err)
	})
}

func TestSlackPost_InvalidURL(t *testing.T) {
	t.Run("invalid URL returns error", func(t *testing.T) {
		payload := []byte(`{"text": "test"}`)
	err := SlackPost(payload, "://invalid-url", false)
		assert.Error(t, err)
	})

	t.Run("unreachable host returns error", func(t *testing.T) {
		payload := []byte(`{"text": "test"}`)
		err := SlackPost(payload, "http://localhost:99999/nonexistent", false)

		assert.Error(t, err)
	})
}

func TestSlackPost_Payload(t *testing.T) {
	t.Run("sends correct payload", func(t *testing.T) {
		expectedPayload := []byte(`{"text": "Hello, World!"}`)
		var receivedPayload []byte

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			receivedPayload = body
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		err := SlackPost(expectedPayload, server.URL, false)

		assert.NoError(t, err)
		assert.Equal(t, expectedPayload, receivedPayload)
	})

	t.Run("handles empty payload", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		err := SlackPost([]byte{}, server.URL, false)

		assert.NoError(t, err)
	})
}

func TestSlackPost_Headers(t *testing.T) {
	t.Run("sets required headers", func(t *testing.T) {
		var headers http.Header

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers = r.Header
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		payload := []byte(`{"text": "test"}`)
		err := SlackPost(payload, server.URL, false)

		assert.NoError(t, err)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "no-cache", headers.Get("Cache-Control"))
	})
}

func TestSlackPost_Method(t *testing.T) {
	t.Run("uses POST method", func(t *testing.T) {
		var method string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method = r.Method
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		payload := []byte(`{"text": "test"}`)
		err := SlackPost(payload, server.URL, false)

		assert.NoError(t, err)
		assert.Equal(t, "POST", method)
	})
}
