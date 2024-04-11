package service

import (
	"context"
	"fmt"
)

type Sender struct {
	http HttpClient
}

func NewSender(httpClient HttpClient) *Sender {
	return &Sender{
		http: httpClient,
	}
}

func (s *Sender) SendRequestsToURLs(ctx context.Context, urls []string) (map[string]int, error) {
	const errTrace = "Sender.SendRequestsToURLs"

	if len(urls) == 0 {
		return make(map[string]int), nil
	}

	res, err := s.http.SendBatch(ctx, urls)
	if err != nil {
		return nil, fmt.Errorf("%s: requests sending error: %w", errTrace, err)
	}

	return res, nil
}
