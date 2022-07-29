package mapper

import (
	"fmt"
	"github.com/MichalPolinkiewicz/roche/model"
	"github.com/MichalPolinkiewicz/roche/pkg/service"
)

type PingRequestMapper struct{}

func (m *PingRequestMapper) Translate(request *model.PingRequest) (*service.PingRequest, error) {
	if request == nil {
		return nil, fmt.Errorf("request can't be nil")
	}
	return &service.PingRequest{
		Message: request.Message,
	}, nil
}

type PingResponseMapper struct{}

func (m *PingResponseMapper) Translate(response *service.PingResponse) (*model.PingResponse, error) {
	if response == nil {
		return nil, fmt.Errorf("request can't be nil")
	}

	return &model.PingResponse{
		Echo:      response.Echo,
		Timestamp: response.Timestamp,
		Env:       response.Env,
		Version:   response.Version,
	}, nil
}
