package service

// PingDecorator holds common, extra response params
type PingDecorator struct {
	AppVersion string
	Env        string
}

func NewPingDecorator(appVersion, env string) PingDecorator {
	return PingDecorator{AppVersion: appVersion, Env: env}
}
