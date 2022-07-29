package rest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MichalPolinkiewicz/roche/pkg/mapper"
	"github.com/MichalPolinkiewicz/roche/pkg/service"
)

func TestPingHandler_Handle(t *testing.T) {
	type fields struct {
		pingClient     PingClient
		requestMapper  RequestMapper
		responseMapper ResponseMapper
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
		expectedBody string
	}{
		{
			name: "unexpected request method",
			fields: fields{
				pingClient: &dummyPingClientNilResponse{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "localhost:8000", http.NoBody),
			},
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name: "invalid request body",
			fields: fields{
				pingClient: &dummyPingClientNilResponse{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "localhost:8000", strings.NewReader("xxxxx")),
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "nil response from client",
			fields: fields{
				pingClient:     &dummyPingClientNilResponse{},
				requestMapper:  &mapper.PingRequestOneToOneMapper{},
				responseMapper: &mapper.PingResponseOneToOneMapper{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "localhost:8000", strings.NewReader("{\"message\": \"bar1\"}")),
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "request with message and extra fields",
			fields: fields{
				pingClient:     &dummyPingClient200{},
				requestMapper:  &mapper.PingRequestOneToOneMapper{},
				responseMapper: &mapper.PingResponseOneToOneMapper{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "localhost:8000", strings.NewReader("{\"value\": \"xx\", \"message\": \"hello\"}")),
			},
			expectedCode: http.StatusOK,
			expectedBody: "{\"echo\":\"xxx\",\"env\":\"dev\",\"version\":\"1\"}",
		},
		{
			name: "200",
			fields: fields{
				pingClient:     &dummyPingClient200{},
				requestMapper:  &mapper.PingRequestOneToOneMapper{},
				responseMapper: &mapper.PingResponseOneToOneMapper{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "localhost:8000", strings.NewReader("{\"message\": \"bar1\"}")),
			},
			expectedCode: http.StatusOK,
			expectedBody: "{\"echo\":\"xxx\",\"env\":\"dev\",\"version\":\"1\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pingHandler := PingHandler{
				pingClient:     tt.fields.pingClient,
				requestMapper:  tt.fields.requestMapper,
				responseMapper: tt.fields.responseMapper,
			}

			pingHandler.Handle(tt.args.w, tt.args.r)
			require.Equal(t, tt.expectedCode, tt.args.w.Result().StatusCode)

			if tt.expectedBody != "" {
				gotBody, err := io.ReadAll(tt.args.w.Result().Body)
				require.Nil(t, err)
				require.Equal(t, tt.expectedBody, string(gotBody))
			}
		})
	}
}

type dummyPingClientNilResponse struct{}

func (d *dummyPingClientNilResponse) Request(ctx context.Context, r *service.PingRequest) (*service.PingResponse, error) {
	return nil, nil
}

type dummyPingClient200 struct{}

func (d *dummyPingClient200) Request(ctx context.Context, r *service.PingRequest) (*service.PingResponse, error) {
	return &service.PingResponse{
		Echo:    "xxx",
		Env:     "dev",
		Version: "1",
	}, nil
}
