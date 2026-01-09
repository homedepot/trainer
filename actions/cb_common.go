package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"context"
	"errors"
	"fmt"
	"github.com/homedepot/trainer/security"
	"github.com/homedepot/trainer/structs/plan"
	"github.com/juju/loggo"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type cbstate struct {
	inprogress bool
	output     chan *ExecuteResult
	aborted    chan *bool
}

var currcb = cbstate{}

func DoCallback(a ArgStruct, p *plan.Plan, ctx context.Context) (r ExecuteResult) {
	logger := loggo.GetLogger("default")
	logger.Tracef("Executing callback action")

	r.Complete = true

	iargs, err := ParseTemplate(p, a.Args)
	if err != nil {
		logger.Warningf("Parsing action template failed, aborting test")
		r.Err = err
		return
	}
	a.Args = iargs

	// Get the callback method from the Action.
	// If method is not POST or GET, exit function.
	method, err := a.GetArg("method", reflect.TypeOf(""), false)
	if err != nil {
		r.Err = err
		return
	}
	methodstr := ""

	if method == nil {
		methodstr = method.(string)
	}

	if methodstr == "" {
		methodstr = "GET"
	}

	// TODO: should accept more methods.
	if methodstr != "POST" && methodstr != "GET" {
		r.Err = errors.New(fmt.Sprintf("invalid method %s", methodstr))
		return r
	}

	url, err := a.GetArg("url", reflect.TypeOf(""), true)
	if err != nil {
		r.Err = err
		return r
	}
	urlstr := url.(string)

	if urlstr == "" {
		r.Err = errors.New("callback url specified but empty")
		return r
	}

	// Get payload content type.
	pct, err := a.GetArg("url", reflect.TypeOf(""), false)
	if err != nil {
		r.Err = err
		return r
	}

	pctstr := "text/plain"
	if pct != nil && pct.(string) != "" {
		pctstr = pct.(string)
	}

	authstr := ""
	auth, ok := a.Args["auth_header"]
	if ok {
		if auth != nil && auth.(string) != "" {
			authstr = auth.(string)
		}
	}

	logger.Debugf("Executing callback to %s", urlstr)

	// Get payload body.
	payloadbody := ""

	payload, err := a.GetArg("payload", reflect.TypeOf(""), false)
	if err != nil {
		logger.Warningf("invalid payload variable: %s", err)
		r.Err = err
		return r
	}

	if payload != nil && payload.(string) != "" {
		// Validate payload file path to prevent path traversal
		if err := security.ValidatePath(payload.(string), ""); err != nil {
			logger.Warningf("Payload file path validation failed: %s", err)
			return
		}
		
		out, err := os.ReadFile(payload.(string))
		if err != nil {
			logger.Warningf("Couldn't read payload file: %s", payload.(string))
			return
		}
		payloadbody = string(out)
	} else {
		logger.Infof("callback does not have a payload body")
	}

	var resp *http.Response
	var req *http.Request
	if method == "GET" {
		logger.Tracef("Sending GET")
		req, err = http.NewRequest("GET", urlstr, strings.NewReader(""))
	} else if method == "POST" {
		logger.Tracef("Sending POST")
		req, err = http.NewRequest("POST", urlstr, strings.NewReader(ParseStringTemplate(p, payloadbody)))
	} else {
		panic("We checked the methods before, but one seems to have slipped through.  FIXME.")
	}
	if err != nil {
		logger.Warningf("NewRequest failed: %s", err.Error())
		r.Err = err
		return
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", pctstr)
	if authstr != "" {
		req.Header.Set("Authorization", authstr)
	}
	hdrs, ok := a.Args["headers"]
	if ok {
		headers := hdrs.(map[interface{}]interface{})
		for k, v := range headers {
			kstr, ok := k.(string)
			if !ok {
				logger.Warningf("header key %s is not a string", k)
				r.Err = errors.New(fmt.Sprintf("header key %s is not a string", k))
				return
			}
			str, ok := v.(string)
			if !ok {
				logger.Warningf("header value %s is not a string", k)
				r.Err = errors.New(fmt.Sprintf("header key %s is not a string", k))
				return
			}
			req.Header.Set(kstr, str)
		}
	}
	client := http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		logger.Warningf("Couldn't execute callback: %s", err.Error())
		r.Err = err
		return
	}
	defer resp.Body.Close()
	var ignorefailure bool
	ign, ok := a.Args["ignore_failure"]
	if ok {
		ignorefailure = ign.(bool)
	}
	if ignorefailure == true {
		logger.Tracef("ignore_failure set, ignoring response code")
	} else if resp.StatusCode != 200 {
		r.Err = fmt.Errorf("callback did not succeed (code %v)", resp.StatusCode)
		return
	}

	rs, err := io.ReadAll(resp.Body)
	if err != nil {
		r.Err = err
		return
	}
	var i interface{}
	rtype, err := a.GetArg("response_type", reflect.TypeOf(""), false)
	if err != nil {
		logger.Warningf("couldn't get response_type variable: %s", err)
		r.Err = err
		return
	}
	if rtype != nil && rtype.(string) == "json" {
		i, err = LoadJSON(string(rs))
		if err != nil {
			r.Err = err
			return
		}
	} else if rtype != nil && rtype.(string) == "yaml" {
		i, err = LoadYAML(string(rs))
		if err != nil {
			r.Err = err
			return
		}
	} else if rtype != nil && rtype.(string) == "string" {
		// do nothing, response_map won't be set, but response will.
		// don't set response_map if you want to use a string.
	} else {
		r.Err = errors.New(fmt.Sprintf("unknown response type: %s", rtype.(string)))
		return
	}
	imap, imapok := i.(map[string]interface{})
	if imapok {
		logger.Tracef("imap: %+v", imap)
		saveresponseasmap, ok := a.Args["save_response_map"]
		if ok && saveresponseasmap.(string) != "" {
			_, ok1 := p.State.Variables[saveresponseasmap.(string)]
			if !ok1 {
				p.State.Variables[saveresponseasmap.(string)] = make(map[string]interface{}, 0)
			}
			logger.Tracef("saving response as map %s", saveresponseasmap)
			p.State.Variables[saveresponseasmap.(string)] = imap
		}
	}
	saveresponse, ok := a.Args["save_response"]
	if ok && saveresponse.(string) != "" {
		logger.Tracef("saving response as %s", saveresponse)
		err := p.State.SetVariable(saveresponse.(string), string(rs))
		if err != nil {
			r.Err = err
			return
		}
	}
	save, ok := a.Args["save"]
	if ok {
		for _, v := range save.([]interface{}) {
			logger.Tracef("Attempting save into %s", v)
			q, ok := imap[v.(string)]
			if !ok {
				logger.Warningf("Could not save %s from json: not present", v)
				continue
			}
			logger.Tracef("saving %s into %s", q.(string), v.(string))
			p.State.Variables[v.(string)] = q.(string)
		}
	}
	r.Success = true
	return
}
