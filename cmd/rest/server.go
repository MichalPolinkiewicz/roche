package rest

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strings"
)

type RestServer struct {
	port      string
	endpoints []*Endpoint
	srv       *http.Server
}

func NewRestServer(port string, endpoints []*Endpoint) (*RestServer, error) {
	if strings.TrimSpace(port) == "" {
		return nil, fmt.Errorf("port can't be empty")
	}
	return &RestServer{port: port, endpoints: endpoints}, nil
}

func (s *RestServer) Run(ctx context.Context) func() error {
	return func() error {
		s.srv = &http.Server{
			Addr: "localhost:" + s.port,
		}

		for _, endpoint := range s.endpoints {
			http.HandleFunc(endpoint.path, endpoint.handler)
		}

		group, ctx := errgroup.WithContext(ctx)
		group.Go(func() error { return s.start() })
		group.Go(func() error { return s.shutdownWatcher(ctx) })
		return nil
	}
}

func (s *RestServer) start() error {
	log.Println("restServer started:", s.port)
	return s.srv.ListenAndServe()
}

func (s *RestServer) shutdownWatcher(ctx context.Context) error {
	<-ctx.Done()
	log.Println("restServer closed")
	return s.srv.Shutdown(ctx)
}
