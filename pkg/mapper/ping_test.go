package mapper

import (
	"github.com/MichalPolinkiewicz/roche/model"
	"github.com/MichalPolinkiewicz/roche/pkg/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPingRequestMapper_Translate(t *testing.T) {
	type args struct {
		request *model.PingRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *service.PingRequest
		wantErr bool
	}{
		{
			name: "nil request - should return nil, error",
			args: args{
				request: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ok - should return request, error",
			args: args{
				request: &model.PingRequest{
					Message: "msg",
				},
			},
			want: &service.PingRequest{
				Message: "msg",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PingRequestMapper{}
			got, err := m.Translate(tt.args.request)

			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestPingResponseMapper_Translate(t *testing.T) {
	type args struct {
		response *service.PingResponse
	}
	tests := []struct {
		name    string
		args    args
		want    *model.PingResponse
		wantErr bool
	}{
		{
			name: "nil response - should return nil && error",
			args: args{
				response: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "filled reponse - should return response && error",
			args: args{
				response: &service.PingResponse{
					Echo:    "x",
					Env:     "dev",
					Version: "1",
				},
			},
			want: &model.PingResponse{
				Echo:    "x",
				Env:     "dev",
				Version: "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PingResponseMapper{}
			got, err := m.Translate(tt.args.response)

			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}
