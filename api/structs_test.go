package api

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
	"github.com/stretchr/testify/assert"
)

func TestHTTPReturnStruct_WriteOutput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		returnStruct   HTTPReturnStruct
		wantStatusCode int
		wantMessage    string
		wantError      bool
	}{
		{
			name: "successful response with 200",
			returnStruct: HTTPReturnStruct{
				Message:    "Operation completed successfully",
				Error:      false,
				ReturnCode: http.StatusOK,
			},
			wantStatusCode: http.StatusOK,
			wantMessage:    "Operation completed successfully",
			wantError:      false,
		},
		{
			name: "error response with 400",
			returnStruct: HTTPReturnStruct{
				Message:    "Bad request: invalid input",
				Error:      true,
				ReturnCode: http.StatusBadRequest,
			},
			wantStatusCode: http.StatusBadRequest,
			wantMessage:    "Bad request: invalid input",
			wantError:      true,
		},
		{
			name: "error response with 500",
			returnStruct: HTTPReturnStruct{
				Message:    "Internal server error occurred",
				Error:      true,
				ReturnCode: http.StatusInternalServerError,
			},
			wantStatusCode: http.StatusInternalServerError,
			wantMessage:    "Internal server error occurred",
			wantError:      true,
		},
		{
			name: "not found response with 404",
			returnStruct: HTTPReturnStruct{
				Message:    "Resource not found",
				Error:      true,
				ReturnCode: http.StatusNotFound,
			},
			wantStatusCode: http.StatusNotFound,
			wantMessage:    "Resource not found",
			wantError:      true,
		},
		{
			name: "created response with 201",
			returnStruct: HTTPReturnStruct{
				Message:    "Resource created",
				Error:      false,
				ReturnCode: http.StatusCreated,
			},
			wantStatusCode: http.StatusCreated,
			wantMessage:    "Resource created",
			wantError:      false,
		},
		{
			name: "empty message",
			returnStruct: HTTPReturnStruct{
				Message:    "",
				Error:      false,
				ReturnCode: http.StatusOK,
			},
			wantStatusCode: http.StatusOK,
			wantMessage:    "",
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			tt.returnStruct.WriteOutput(c)
			
			assert.Equal(t, tt.wantStatusCode, w.Code, "Status code mismatch")
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Content-Type header should be application/json")
			
			var response HTTPReturnStruct
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")
			
			assert.Equal(t, tt.wantMessage, response.Message, "Message mismatch")
			assert.Equal(t, tt.wantError, response.Error, "Error flag mismatch")
		})
	}
}

func TestHTTPReturnStruct_WriteOutput_JSONMarshaling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("verify ReturnCode is not included in JSON output", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		returnStruct := HTTPReturnStruct{
			Message:    "test message",
			Error:      false,
			ReturnCode: http.StatusOK,
		}
		
		returnStruct.WriteOutput(c)
		
		var jsonMap map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &jsonMap)
		assert.NoError(t, err)
		
		_, hasReturnCode := jsonMap["ReturnCode"]
		assert.False(t, hasReturnCode, "ReturnCode should not be in JSON output due to json:\"-\" tag")
		
		assert.Contains(t, jsonMap, "message")
		assert.Contains(t, jsonMap, "error")
	})
}

func TestHTTPReturnStruct_WriteOutput_SpecialCharacters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "message with quotes",
			message: `Message with "quotes" inside`,
		},
		{
			name:    "message with newlines",
			message: "Message with\nnewlines",
		},
		{
			name:    "message with unicode",
			message: "Message with unicode: 你好世界",
		},
		{
			name:    "message with special chars",
			message: "Message with special chars: <>&'\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			
			returnStruct := HTTPReturnStruct{
				Message:    tt.message,
				Error:      false,
				ReturnCode: http.StatusOK,
			}
			
			returnStruct.WriteOutput(c)
			
			var response HTTPReturnStruct
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Should handle special characters in JSON")
			assert.Equal(t, tt.message, response.Message, "Message should be preserved correctly")
		})
	}
}
