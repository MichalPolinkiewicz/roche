package grpc

import (
	"context"
	proto "github.com/MichalPolinkiewicz/roche/model/proto"
	"time"
)

type GrpcServer struct {
	proto.PingServiceServer
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (s *GrpcServer) Ping(context.Context, *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{
		Echo:      "dummy",
		Timestamp: time.Now().String(),
		Env:       "dummy",
		Version:   "1",
	}, nil
}
