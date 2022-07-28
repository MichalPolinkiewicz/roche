package postman

import (
	"fmt"
	"strings"
	"time"
)

type PostmanServiceConfig struct {
	endpoint       string
	requestTimeout time.Duration
	requestParams  map[string]struct{}
}

func NewPostmanServiceConfig(endpoint string, requestTimeout time.Duration, requestParams []string) (*PostmanServiceConfig, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, fmt.Errorf("no endpoint provided")
	}
	if requestTimeout == 0 {
		return nil, fmt.Errorf("timeout can't be 0")
	}
	return &PostmanServiceConfig{
		endpoint:       endpoint,
		requestTimeout: requestTimeout,
		requestParams:  requestParamsMapper(requestParams),
	}, nil
}

func requestParamsMapper(params []string) map[string]struct{} {
	results := make(map[string]struct{})
	for _, param := range params {
		results[param] = struct{}{}
	}
	return results
}
