package service

import (
	"context"
	"fmt"
	"time"
)

type pingable interface {
	Get(ctx context.Context, message string) ([]byte, error)
}

type pingService struct {
	timeout          time.Duration
	reponseDecorator PingDecorator
	pingable         pingable
}

func NewPingService(pingable pingable, timeout time.Duration, decorator PingDecorator) (*pingService, error) {
	if pingable == nil {
		return nil, fmt.Errorf("pingable can't be nil")
	}
	if timeout == 0 {
		return nil, fmt.Errorf("timeout can't be 0")
	}
	return &pingService{
		timeout:          timeout,
		reponseDecorator: decorator,
		pingable:         pingable,
	}, nil
}

func (s *pingService) Request(ctx context.Context, r *PingRequest) (*PingResponse, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, s.timeout)
	defer cancelFunc()

	echo, err := s.pingable.Get(ctx, r.Message)
	if err != nil {
		return nil, err
	}

	return &PingResponse{
		Echo:      string(echo),
		Timestamp: time.Now().Unix(),
		Env:       s.reponseDecorator.Env,
		Version:   s.reponseDecorator.AppVersion,
	}, nil
}
