package service

import "context"

type HttpClient interface {
	SendBatch(ctx context.Context, urls []string) (map[string]int, error)
}
