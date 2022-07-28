package rest

import (
	"fmt"
	"strings"
)

type RestServerConfig struct {
	Port string
}

func NewRestServerConfig(port string) (*RestServerConfig, error) {
	if strings.TrimSpace(port) == "" {
		return nil, fmt.Errorf("no valid port provided")
	}
	return &RestServerConfig{Port: port}, nil
}
