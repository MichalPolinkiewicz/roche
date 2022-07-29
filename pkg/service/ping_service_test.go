package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestNewPingService(t *testing.T) {
	type args struct {
		pingable  pingable
		timeout   time.Duration
		decorator PingDecorator
	}
	tests := []struct {
		name    string
		args    args
		want    *pingService
		wantErr bool
	}{
		{
			name: "no client, should return nil && error",
			args: args{
				pingable: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "0 timeout - should return nil && error",
			args: args{
				pingable: &dummyPingableMock{},
				timeout:  0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid params - should return service && nil",
			args: args{
				pingable: &dummyPingableMock{},
				timeout:  time.Second,
			},
			want: &pingService{
				timeout:          time.Second,
				reponseDecorator: PingDecorator{},
				pingable:         &dummyPingableMock{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPingService(tt.args.pingable, tt.args.timeout, tt.args.decorator)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPingService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPingService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type dummyPingableMock struct{}

func (p *dummyPingableMock) Get(ctx context.Context, message string) ([]byte, error) {
	return nil, nil
}

func TestPingService_RequestClientErrorResponse(t *testing.T) {
	s := &pingService{
		timeout:  time.Second,
		pingable: &dummyErrorPingableMock{},
	}
	resp, err := s.Request(context.Background(), &PingRequest{})
	require.NotNil(t, err)
	require.Nil(t, resp)
}

type dummyErrorPingableMock struct{}

func (p *dummyErrorPingableMock) Get(ctx context.Context, message string) ([]byte, error) {
	return nil, fmt.Errorf("dummy error")
}

func TestPingService_ClientTimeout(t *testing.T) {
	s := &pingService{
		timeout:  time.Millisecond,
		pingable: &dummyTimeoutPingableMock{},
	}
	resp, err := s.Request(context.Background(), &PingRequest{})
	require.Equal(t, context.DeadlineExceeded, err)
	require.Nil(t, resp)
}

type dummyTimeoutPingableMock struct{}

func (p *dummyTimeoutPingableMock) Get(ctx context.Context, message string) ([]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
			return nil, nil
		}
	}
}

func TestPingService_RequestClient200(t *testing.T) {
	s := &pingService{
		timeout:  time.Second,
		pingable: &dummy200PingableMock{},
		reponseDecorator: PingDecorator{
			AppVersion: "1",
			Env:        "test",
		},
	}
	resp, err := s.Request(context.Background(), &PingRequest{})

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "1", resp.Version)
	require.Equal(t, "test", resp.Env)
	require.Equal(t, "xxx", resp.Echo)
}

type dummy200PingableMock struct{}

func (p *dummy200PingableMock) Get(ctx context.Context, message string) ([]byte, error) {
	return []byte("xxx"), nil
}
