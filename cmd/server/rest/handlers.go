package rest

import (
	"context"
	"fmt"
	proto "github.com/MichalPolinkiewicz/roche/model"
	"github.com/golang/protobuf/jsonpb"
	"net/http"
	"time"
)

type Callable interface {
	Call(ctx context.Context, requestParams map[string]string) ([]byte, error)
}

type MessageHandler struct {
	callable Callable
}

func NewMessageHandler(callable Callable) (*MessageHandler, error) {
	if callable == nil {
		return nil, fmt.Errorf("callable can't be nil")
	}
	return &MessageHandler{
		callable: callable,
	}, nil
}

func (h *MessageHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var pingRequest proto.PingRequest
	if err := jsonpb.Unmarshal(r.Body, &pingRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := h.callable.Call(context.Background(), map[string]string{"message": pingRequest.Message})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pingResponse := proto.PingResponse{
		Echo:      string(body),
		Timestamp: time.Now().String(),
		Env:       "",
		Version:   "",
	}

	m := jsonpb.Marshaler{}
	if err := m.Marshal(w, &pingResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
