package mapper

import (
	"fmt"

	"github.com/MichalPolinkiewicz/roche/model"
	"github.com/MichalPolinkiewicz/roche/pkg/service"
)

type PingRequestOneToOneMapper struct{}

func (m *PingRequestOneToOneMapper) Translate(request *model.PingRequest) (*service.PingRequest, error) {
	if request == nil {
		return nil, fmt.Errorf("request can't be nil")
	}
	return &service.PingRequest{
		Message: request.Message,
	}, nil
}

type PingResponseOneToOneMapper struct{}

func (m *PingResponseOneToOneMapper) Translate(response *service.PingResponse) (*model.PingResponse, error) {
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
