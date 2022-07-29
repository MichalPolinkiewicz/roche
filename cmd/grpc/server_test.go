package grpc

import (
	"context"
	"github.com/MichalPolinkiewicz/roche/model"
	"github.com/MichalPolinkiewicz/roche/pkg/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func TestNewGrpcServer(t *testing.T) {
	type args struct {
		port       string
		pingClient PingClient
	}
	tests := []struct {
		name    string
		args    args
		want    *GrpcServer
		wantErr bool
	}{
		{
			name: "empty port - should return nil && error",
			args: args{
				port: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no client - should return nil && error",
			args: args{
				port:       "",
				pingClient: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no client - should return nil && error",
			args: args{
				port:       "3000",
				pingClient: &dummyPingClientMock{},
			},
			want: &GrpcServer{
				port:       "3000",
				pingClient: &dummyPingClientMock{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGrpcServer(tt.args.port, tt.args.pingClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGrpcServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGrpcServer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type dummyPingClientMock struct{}

func (m *dummyPingClientMock) Request(ctx context.Context, r *service.PingRequest) (*service.PingResponse, error) {
	return nil, nil
}

func TestGrpcServer_Ping200(t *testing.T) {
	port := "30069"
	srv, err := NewGrpcServer(port, &dummyPingClientMock200{})
	require.Nil(t, err)
	require.NotNil(t, srv)

	ctx, cancel := context.WithCancel(context.Background())

	go func() { srv.Run(ctx)() }()
	defer cancel()

	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	defer conn.Close()

	client := NewPingServiceClient(conn)
	got, err := client.Ping(context.Background(), &model.PingRequest{Message: "aaaa"})

	require.Nil(t, err)
	require.NotNil(t, got)
	require.True(t, proto.Equal(&model.PingResponse{
		Echo:      "grpc test",
		Timestamp: "",
		Env:       "dev",
		Version:   "2",
	}, got))
}

type dummyPingClientMock200 struct{}

func (m *dummyPingClientMock200) Request(ctx context.Context, r *service.PingRequest) (*service.PingResponse, error) {
	return &service.PingResponse{
		Echo:      "grpc test",
		Timestamp: "",
		Env:       "dev",
		Version:   "2",
	}, nil
}

func TestGrpcServer_PingNilRequest(t *testing.T) {
	port := "30069"
	srv, err := NewGrpcServer(port, &dummyPingClientMock{})
	require.Nil(t, err)
	require.NotNil(t, srv)

	ctx, cancel := context.WithCancel(context.Background())

	go func() { srv.Run(ctx)() }()
	defer cancel()

	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	defer conn.Close()

	client := NewPingServiceClient(conn)
	got, err := client.Ping(context.Background(), nil)

	require.NotNil(t, err)
	require.Nil(t, got)
}
