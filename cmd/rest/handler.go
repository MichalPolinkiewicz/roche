package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MichalPolinkiewicz/roche/model"
	"github.com/MichalPolinkiewicz/roche/pkg/service"
	"github.com/golang/protobuf/jsonpb"
)

type PingClient interface {
	Request(ctx context.Context, r *service.PingRequest) (*service.PingResponse, error)
}

type RequestMapper interface {
	Translate(request *model.PingRequest) (*service.PingRequest, error)
}

type ResponseMapper interface {
	Translate(response *service.PingResponse) (*model.PingResponse, error)
}

type PingHandler struct {
	pingClient     PingClient
	requestMapper  RequestMapper
	responseMapper ResponseMapper
}

func NewPingHandler(pingClient PingClient, reqMapper RequestMapper, respMapper ResponseMapper) (*PingHandler, error) {
	if pingClient == nil {
		return nil, fmt.Errorf("pingable can't be nil")
	}
	if reqMapper == nil || respMapper == nil {
		return nil, fmt.Errorf("request and response mappers can't be nil")
	}
	return &PingHandler{
		pingClient:     pingClient,
		responseMapper: respMapper,
		requestMapper:  reqMapper,
	}, nil
}

func (h *PingHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var pingRequest model.PingRequest
	unmarshaller := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err := unmarshaller.Unmarshal(r.Body, &pingRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	serviceRequest, err := h.requestMapper.Translate(&pingRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	serviceResponse, err := h.pingClient.Request(r.Context(), serviceRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pingResponse, err := h.responseMapper.Translate(serviceResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := (&jsonpb.Marshaler{}).Marshal(w, pingResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
