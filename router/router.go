package router

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homedepot/trainer/api"
	"github.com/homedepot/trainer/api/apiv1"
	"github.com/homedepot/trainer/cli"
	"github.com/homedepot/trainer/config"
	"github.com/homedepot/trainer/handler"
	"github.com/juju/loggo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"sync"
)

type Router struct {
	Router         *gin.Engine
	NoRouteHandler handler.Handler
}

func StartRouter(o cli.Options, c *config.Config, wg *sync.WaitGroup, h *handler.Handler) {
	logger := loggo.GetLogger("default")
	level, _ := loggo.ParseLevel(o.LogLevel)
	logger.SetLogLevel(level)
	go func(o cli.Options) {
		router := Router{}
		router.CreateRouter(o, c, logger, nil, h)
		router.Run(o, logger)
		wg.Done()
		return
	}(o)
}

func Redirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme := c.GetHeader("X-Forwarded-Proto")

		if scheme == "http" {
			// Write an error and stop the handler chain
			c.Abort()
			c.Redirect(http.StatusMovedPermanently, "https://"+c.Request.Host+c.Request.URL.String())
		}
	}
}

func (r *Router) CreateRouter(o cli.Options, c *config.Config, logger loggo.Logger, e *gin.Engine, h *handler.Handler) {
	if e == nil {
		e = gin.Default()
	}

	e.Use(Redirect())
	e.GET("/health-check", api.HealthCheck)
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))
	e.OPTIONS("/cloudfoundryapplication", api.CloudFoundry)

	e.NoRoute(r.AddHandler(h))

	version1 := apiv1.V1{}
	version1.AddHandlers(o, c, e, h)
	r.Router = e
}

func (r *Router) AddHandler(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.Add(c)
	}
}

func (r *Router) Run(o cli.Options, logger loggo.Logger) {
	err := r.Router.Run(fmt.Sprintf("%s:%s", o.APIListenHost, o.APIListenPort))
	if err != nil {
		logger.Criticalf("Couldn't start GIN")
		os.Exit(1)
	}
}
