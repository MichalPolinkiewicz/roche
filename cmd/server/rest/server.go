package rest

import "net/http"

type RestServer struct {
	Config    *RestServerConfig
	Endpoints []*Endpoint
}

func NewRestServer(config *RestServerConfig, endpoints []*Endpoint) *RestServer {
	return &RestServer{Config: config, Endpoints: endpoints}
}

func (s *RestServer) Run() error {
	for _, endpoint := range s.Endpoints {
		http.HandleFunc(endpoint.path, endpoint.handler)
	}
	if err := http.ListenAndServe(":"+s.Config.Port, nil); err != nil {
		return err
	}
	return nil
}
