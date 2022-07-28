package grpc

import (
	"context"
	"github.com/MichalPolinkiewicz/roche/model"
	"time"
)

type GrpcServer struct {
	PingServiceServer
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (s *GrpcServer) Ping(context.Context, *model.PingRequest) (*model.PingResponse, error) {
	return &model.PingResponse{
		Echo:      "dummy",
		Timestamp: time.Now().String(),
		Env:       "dummy",
		Version:   "1",
	}, nil
}
