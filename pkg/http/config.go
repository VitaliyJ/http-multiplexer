package http

import "time"

const defaultMaxConcurrentRequests = 4

const defaultRequestTimeout = time.Second

type Config struct {
	MaxConcurrentRequests int
	RequestTimeout        time.Duration
	AllowedURLs           []string
	BlockedURLs           []string
	AllowAll              bool
}
