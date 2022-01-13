package apiv1

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/homedepot/trainer/api"
	"github.com/homedepot/trainer/cli"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/handler"
	"github.com/juju/loggo"
	"gopkg.in/yaml.v2"
)

type V1 struct {
	api.Version
}

func (v *V1) AddHandlers(o cli.Options, c *config.Config, g *gin.Engine, h *handler.Handler) {

	logger := loggo.GetLogger("default")
	level, _ := loggo.ParseLevel(o.LogLevel)
	logger.SetLogLevel(level)
	group := g.Group("/capi/v1", gin.BasicAuth(gin.Accounts{o.APIAuthUsername: o.APIAuthPass}))
	{
		group.POST("/launch/:plan", v.Launch(c))
		group.POST("/remove", v.Remove(c, h))
		group.POST("/status", v.Status(c))
		group.POST("/config", v.ConfigAPI(c))

	}
}

// Plan the endpoint for setting plan according to request from user.
func (v *V1) Launch(cfg *config.Config) func(*gin.Context) {
	return func(c *gin.Context) {

		p := c.Param("plan")

		err := handler.LaunchTest(cfg, p)

		if err != nil {
			resp := api.HTTPReturnStruct{
				Message:    err.Error(),
				Error:      true,
				ReturnCode: 400,
			}
			resp.WriteOutput(c)
			return
		}

		resp := api.HTTPReturnStruct{
			Message:    "Plan changed successfully",
			Error:      false,
			ReturnCode: 200,
		}
		resp.WriteOutput(c)
		return
	}
}

// ConfigAPI the endpoint for setting desired testing configuration.
func (v *V1) ConfigAPI(cfg *config.Config) func(*gin.Context) {
	return func(c *gin.Context) {

		configstr, err := yaml.Marshal(cfg)
		if err != nil {
			resp := api.HTTPReturnStruct{
				Message:    err.Error(),
				Error:      true,
				ReturnCode: 400,
			}
			resp.WriteOutput(c)
			return
		}
		c.Writer.WriteHeader(200)
		_, _ = c.Writer.Write(configstr)
	}
}

// Status returns Status object reflecting current plan,
// transaction, states, and runner switch flag.
func (v *V1) Status(cfg *config.Config) func(*gin.Context) {
	return func(c *gin.Context) {
		logger := loggo.GetLogger("default")

		s := handler.GetPlan()
		if s == nil {
			resp := api.HTTPReturnStruct{
				Error:      false,
				Message:    "no change",
				ReturnCode: 200,
			}
			resp.WriteOutput(c)
			return
		}
		logger.Tracef("Status: %+v", s)

		sout, err := json.Marshal(s.State)
		if err != nil {
			logger.Warningf("Error marshaling JSON: %s", err)
			resp := api.HTTPReturnStruct{
				Message:    err.Error(),
				Error:      true,
				ReturnCode: 400,
			}
			resp.WriteOutput(c)
			return
		}
		c.Writer.WriteHeader(200)
		c.Writer.Write(sout)
	}
}

// Remove removes a test
func (v *V1) Remove(cfg *config.Config, h *handler.Handler) func(*gin.Context) {
	return func(c *gin.Context) {

		h.Reset()
		c.Writer.WriteHeader(200)
		c.Writer.Write([]byte(`{"message": "remove succeeded", "error": false}`))
	}
}
