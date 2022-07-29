package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Callable interface {
	Do(req *http.Request) (*http.Response, error)
}

type PostmanClient struct {
	endpoint string
	callable Callable
}

func NewPostmanClient(endpoint string, callable Callable) (*PostmanClient, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, fmt.Errorf("endpoint can't be nil")
	}
	return &PostmanClient{endpoint: endpoint, callable: callable}, nil
}

func (c *PostmanClient) Get(ctx context.Context, message string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.URL.Query().Add("message", message)

	resp, err := c.callable.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service responded with unexpected code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = resp.Body.Close(); err != nil {
		log.Printf("can't close response body: %v", err)
	}
	return respBody, nil
}
