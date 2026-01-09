package core

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TODO : Need to figure out how to access Slack variables in Vault (not working like in pipeline files)

// SlackPost posts payload to Slack.
// insecureSkipVerify should only be set to true in trusted internal networks.
// In production, this should be false to ensure TLS certificate validation.
func SlackPost(payload []byte, url string, insecureSkipVerify bool) error {

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	proxyURL, err := http.ProxyFromEnvironment(request)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		Proxy:           http.ProxyURL(proxyURL),
	}

	httpClient := &http.Client{
		Timeout:   time.Second * 30,
		Transport: tr,
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Cache-Control", "no-cache")

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	_, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= 300 || response.StatusCode < 200 {
		return fmt.Errorf("slack post failed with status code: %d", response.StatusCode)
	}

	return nil
}
