package http

const defaultMaxRequestsCount int64 = 100

type Config struct {
	Port             string
	MaxRequestsCount int64
}
