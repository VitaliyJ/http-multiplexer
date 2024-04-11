package http

import "context"

type Sender interface {
	SendRequestsToURLs(ctx context.Context, urls []string) (map[string]int, error)
}
