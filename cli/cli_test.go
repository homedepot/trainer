package cli

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOptions(t *testing.T) {
	t.Run("NewOptions returns non-nil Options", func(t *testing.T) {
		o := NewOptions()
		assert.NotNil(t, o)
	})

	t.Run("NewOptions returns empty Options struct", func(t *testing.T) {
		o := NewOptions()
		assert.Equal(t, "", o.APIListenPort)
		assert.Equal(t, "", o.APIListenHost)
		assert.Equal(t, "", o.APIAuthUsername)
		assert.Equal(t, "", o.APIAuthPass)
		assert.Equal(t, "", o.LogLevel)
		assert.Equal(t, "", o.ConfigFile)
		assert.False(t, o.TestMode)
	})
}

func TestOptions_Parse_Defaults(t *testing.T) {
	t.Run("Parse sets default values", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
		}

		o.Parse(args)

		assert.Equal(t, "8080", o.APIListenPort, "Default port should be 8080")
		assert.Equal(t, "localhost", o.APIListenHost, "Default host should be localhost")
		assert.Equal(t, "WARNING", o.LogLevel, "Default log level should be WARNING")
		assert.Equal(t, "config.yml", o.ConfigFile, "Default config file should be config.yml")
		assert.False(t, o.TestMode, "Default test mode should be false")
		assert.Equal(t, "", o.TestURL, "Default test URL should be empty")
	})
}

func TestOptions_Parse_CustomValues(t *testing.T) {
	t.Run("Parse accepts custom port", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
			"--apiport", "9090",
		}

		o.Parse(args)

		assert.Equal(t, "9090", o.APIListenPort)
	})

	t.Run("Parse accepts custom host", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
			"--apihost", "0.0.0.0",
		}

		o.Parse(args)

		assert.Equal(t, "0.0.0.0", o.APIListenHost)
	})

	t.Run("Parse accepts custom log level", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
			"--loglevel", "DEBUG",
		}

		o.Parse(args)

		assert.Equal(t, "DEBUG", o.LogLevel)
	})

	t.Run("Parse accepts custom config file", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
			"--configfile", "custom.yml",
		}

		o.Parse(args)

		assert.Equal(t, "custom.yml", o.ConfigFile)
	})

	t.Run("Parse accepts test mode", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
			"--testmode",
		}

		o.Parse(args)

		assert.True(t, o.TestMode)
	})

	t.Run("Parse accepts test URL", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
			"--testurl", "http://test.example.com",
		}

		o.Parse(args)

		assert.Equal(t, "http://test.example.com", o.TestURL)
	})

	t.Run("Parse accepts bases map", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
			"--bases", "key1=value1",
			"--bases", "key2=value2",
		}

		o.Parse(args)

		assert.NotNil(t, o.Bases)
		assert.Equal(t, "value1", o.Bases["key1"])
		assert.Equal(t, "value2", o.Bases["key2"])
	})
}

func TestOptions_Parse_RequiredFields(t *testing.T) {
	t.Run("Parse succeeds with both required fields", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
		}

		o.Parse(args)

		assert.Equal(t, "testuser", o.APIAuthUsername)
		assert.Equal(t, "testpass", o.APIAuthPass)
	})
}

func TestOptions_Parse_EnvironmentVariables(t *testing.T) {
	t.Run("Parse reads from environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("PORT", "7070")
		os.Setenv("APILISTENHOST", "127.0.0.1")
		os.Setenv("APIAUTHUSERNAME", "envuser")
		os.Setenv("APIAUTHPASSWORD", "envpass")
		os.Setenv("LOGLEVEL", "INFO")
		os.Setenv("CONFIGFILE", "env.yml")
		os.Setenv("TESTMODE", "true")
		os.Setenv("TESTURL", "http://env.example.com")

		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("APILISTENHOST")
			os.Unsetenv("APIAUTHUSERNAME")
			os.Unsetenv("APIAUTHPASSWORD")
			os.Unsetenv("LOGLEVEL")
			os.Unsetenv("CONFIGFILE")
			os.Unsetenv("TESTMODE")
			os.Unsetenv("TESTURL")
		}()

		o := NewOptions()
		args := []string{} // Empty args, should use env vars

		o.Parse(args)

		assert.Equal(t, "7070", o.APIListenPort)
		assert.Equal(t, "127.0.0.1", o.APIListenHost)
		assert.Equal(t, "envuser", o.APIAuthUsername)
		assert.Equal(t, "envpass", o.APIAuthPass)
		assert.Equal(t, "INFO", o.LogLevel)
		assert.Equal(t, "env.yml", o.ConfigFile)
		assert.True(t, o.TestMode)
		assert.Equal(t, "http://env.example.com", o.TestURL)
	})

	t.Run("Parse command line flags override environment variables", func(t *testing.T) {
		os.Setenv("PORT", "7070")
		os.Setenv("APIAUTHUSERNAME", "envuser")
		os.Setenv("APIAUTHPASSWORD", "envpass")

		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("APIAUTHUSERNAME")
			os.Unsetenv("APIAUTHPASSWORD")
		}()

		o := NewOptions()
		args := []string{
			"--apiuser", "flaguser",
			"--apipass", "flagpass",
			"--apiport", "9090",
		}

		o.Parse(args)

		assert.Equal(t, "9090", o.APIListenPort, "Flag should override env var")
		assert.Equal(t, "flaguser", o.APIAuthUsername, "Flag should override env var")
		assert.Equal(t, "flagpass", o.APIAuthPass, "Flag should override env var")
	})
}

func TestOptions_Parse_App(t *testing.T) {
	t.Run("Parse sets App field", func(t *testing.T) {
		o := NewOptions()
		args := []string{
			"--apiuser", "testuser",
			"--apipass", "testpass",
		}

		o.Parse(args)

		assert.NotNil(t, o.App)
	})
}
