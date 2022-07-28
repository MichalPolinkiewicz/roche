package postman

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type PostmanService struct {
	client http.Client
	config *PostmanServiceConfig
}

func NewPostmanService(config *PostmanServiceConfig, client http.Client) (*PostmanService, error) {
	if config == nil {
		return nil, fmt.Errorf("no config provided")
	}
	return &PostmanService{config: config, client: client}, nil
}

func (s *PostmanService) Call(ctx context.Context, requestParams map[string]string) ([]byte, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, s.config.requestTimeout)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.config.endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}

	for k, v := range requestParams {
		if _, ok := s.config.allowedRequestParams[k]; ok {
			req.URL.Query().Add(k, v)
		}
	}

	resp, err := s.client.Do(req)
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
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return respBody, nil
}
