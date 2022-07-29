package service

type PingResponse struct {
	Echo      string
	Timestamp int64
	Env       string
	Version   string
}
