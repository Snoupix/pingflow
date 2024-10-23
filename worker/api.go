package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"
)

const API_BASE_URL = "https://www.dnd5eapi.co"
const API_RESOURCES = "/api"
const CLASSES = "/api/classes"
const SUBCLASSES = "/api/subclasses"

func FetchEndpoint(ctx context.Context, httpclient *http.Client, config WorkConfig) (string, bool) {
	enpoint, ok := ParseEndpoint(config)
	if !ok {
		return "{ \"error\": \"This endpoint is not implemented by the worker\" }", false
	}

	url := API_BASE_URL + enpoint
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Could not build request: %s", err)
	}

	resp, err := httpclient.Do(req)
	if err != nil {
		log.Fatalf("Could not send request: %s", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Could not read response body: %s", err)
	}

	return string(data), true
}

// This will return empty string and false if this specific endpoint parsing isn't implemented
// Otherwise, it returns the parsed endpoint without the base and true
func ParseEndpoint(config WorkConfig) (string, bool) {
	// Potential attack vector; TODO: sanitize end result
    config.endpoint = strings.TrimSuffix(config.endpoint, "/")

	switch config.endpoint {
	case API_RESOURCES, CLASSES, SUBCLASSES:
		return config.endpoint + "/" + config.parameters, true
	}

	return "", false
}
