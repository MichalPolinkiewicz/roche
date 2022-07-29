package config

import (
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Env     string `mapstructure:"env"`
	Version string `mapstructure:"version"`

	GrpcPort string `mapstructure:"grpc_port"`
	RestPort string `mapstructure:"rest_port"`

	PingServiceTimeout        time.Duration `mapstructure:"ping_service_request_timeout"`
	PingServiceClientEndpoint string        `mapstructure:"ping_service_client_endpoint"`
}

func NewAppConfig(filename string) (*AppConfig, error) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var appConfig AppConfig
	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}
