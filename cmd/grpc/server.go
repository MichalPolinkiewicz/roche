package grpc

import (
	"context"
	"fmt"
	"github.com/MichalPolinkiewicz/roche/model"
	"github.com/MichalPolinkiewicz/roche/pkg/mapper"
	"github.com/MichalPolinkiewicz/roche/pkg/service"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

type PingClient interface {
	Request(ctx context.Context, r *service.PingRequest) (*service.PingResponse, error)
}

type GrpcServer struct {
	port       string
	pingClient PingClient
	srv        *grpc.Server
}

func NewGrpcServer(port string, pingClient PingClient) (*GrpcServer, error) {
	if strings.TrimSpace(port) == "" {
		return nil, fmt.Errorf("invalid port provided")
	}
	return &GrpcServer{
		port:       port,
		pingClient: pingClient,
	}, nil
}

func (s *GrpcServer) Run(ctx context.Context) func() error {
	return func() error {
		s.srv = grpc.NewServer()
		listener, err := net.Listen("tcp", ":"+s.port)
		if err != nil {
			return err
		}

		RegisterPingServiceServer(s.srv, s)
		if err != nil {
			return err
		}

		group, ctx := errgroup.WithContext(ctx)
		group.Go(func() error { return s.start(listener) })
		group.Go(func() error { return s.shutdownWatcher(ctx) })
		if err := group.Wait(); err != nil {
			return err
		}
		return nil
	}
}

func (s *GrpcServer) start(listener net.Listener) error {
	log.Println("grpcServer started:", s.port)
	return s.srv.Serve(listener)
}

func (s *GrpcServer) shutdownWatcher(ctx context.Context) error {
	<-ctx.Done()
	s.srv.GracefulStop()
	log.Println("grpcServer closed")
	return nil
}

func (s *GrpcServer) Ping(ctx context.Context, pingRequest *model.PingRequest) (*model.PingResponse, error) {
	requestMapper, responseMapper := mapper.PingRequestMapper{}, mapper.PingResponseMapper{}
	serviceRequest, err := requestMapper.Translate(pingRequest)
	if err != nil {
		return nil, err
	}
	serviceResp, err := s.pingClient.Request(ctx, serviceRequest)
	if err != nil {
		return nil, err
	}
	pingResponse, err := responseMapper.Translate(serviceResp)
	if err != nil {
		return nil, err
	}
	return pingResponse, nil
}

func (s *GrpcServer) mustEmbedUnimplementedPingServiceServer() {}
