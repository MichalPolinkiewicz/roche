package rest

import (
	"fmt"
	"net/http"
	"strings"
)

type Endpoint struct {
	path    string
	method  string
	handler http.HandlerFunc
}

func NewEndpoint(path, method string, handler http.HandlerFunc) (*Endpoint, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("invalid path")
	}
	if method != http.MethodGet && method != http.MethodPost {
		return nil, fmt.Errorf("invalid method: %v", method)
	}
	if handler == nil {
		return nil, fmt.Errorf("handler can't be nil")
	}
	return &Endpoint{
		path:    path,
		method:  method,
		handler: handler,
	}, nil
}
