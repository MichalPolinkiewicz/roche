package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestPostmanClient_GetClientResponseError(t *testing.T) {
	postmanClient, err := NewPostmanClient("dummy/1234", &dummyCallableMock{})
	require.Nil(t, err)
	require.NotNil(t, postmanClient)

	resp, err := postmanClient.Get(context.Background(), "")
	require.Nil(t, resp)
	require.NotNil(t, err)
}

type dummyCallableMock struct{}

func (m *dummyCallableMock) Do(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("err")
}

func TestPostmanClient_GetClientResponse500(t *testing.T) {
	postmanClient, err := NewPostmanClient("dummy/1234", &dummyCallableMock500{})
	require.Nil(t, err)
	require.NotNil(t, postmanClient)

	resp, err := postmanClient.Get(context.Background(), "")
	require.Nil(t, resp)
	require.NotNil(t, err)
}

type dummyCallableMock500 struct{}

func (m *dummyCallableMock500) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusInternalServerError}, nil
}

func TestPostmanClient_GetClientResponse200NoBody(t *testing.T) {
	postmanClient, err := NewPostmanClient("dummy/1234", &dummyCallableMock200NoBody{})
	require.Nil(t, err)
	require.NotNil(t, postmanClient)

	resp, err := postmanClient.Get(context.Background(), "")
	require.Nil(t, err)
	require.NotNil(t, resp)
}

type dummyCallableMock200NoBody struct{}

func (m *dummyCallableMock200NoBody) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
}

func TestPostmanClient_GetClientResponse200WithBody(t *testing.T) {
	postmanClient, err := NewPostmanClient("dummy/1234", &dummyCallableMock200WithBody{})
	require.Nil(t, err)
	require.NotNil(t, postmanClient)

	resp, err := postmanClient.Get(context.Background(), "")
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "xxx", string(resp))
}

type dummyCallableMock200WithBody struct{}

func (m *dummyCallableMock200WithBody) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("xxx")),
		},
		nil
}

func TestPostmanClient_GetClientResponse200JSON(t *testing.T) {
	postmanClient, err := NewPostmanClient("dummy/1234", &dummyCallableMock200JSON{})
	require.Nil(t, err)
	require.NotNil(t, postmanClient)

	resp, err := postmanClient.Get(context.Background(), "")
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "{\"args\": {\"foo1\": \"bar1\",\"foo2\": \"bar2\"}}", string(resp))
}

type dummyCallableMock200JSON struct{}

func (m *dummyCallableMock200JSON) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("{\"args\": {\"foo1\": \"bar1\",\"foo2\": \"bar2\"}}")),
		},
		nil
}
