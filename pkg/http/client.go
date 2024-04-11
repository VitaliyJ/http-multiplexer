package http

import (
	"context"
	"net/http"
	"slices"
	"sync"
)

type Client struct {
	cfg *Config
}

func NewClient(cfg *Config) *Client {
	if cfg.MaxConcurrentRequests == 0 {
		cfg.MaxConcurrentRequests = defaultMaxConcurrentRequests
	}
	if cfg.RequestTimeout == 0 {
		cfg.RequestTimeout = defaultRequestTimeout
	}

	return &Client{
		cfg: cfg,
	}
}

func (c *Client) SendBatch(ctx context.Context, urls []string) (map[string]int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var respErr error
	workersCount := c.cfg.MaxConcurrentRequests
	if len(urls) < c.cfg.MaxConcurrentRequests {
		workersCount = len(urls)
	}
	ch := make(chan string, workersCount)
	wg := sync.WaitGroup{}
	wg.Add(workersCount)

	for i := 0; i < workersCount; i++ {
		go func() {
			defer wg.Done()
			for url := range ch {
				// проверяем, не отменён ли контекст
				if ctx.Err() != nil {
					return
				}
				// отправляем запрос, отменяем конекст в случае ошибки
				if err := c.sendRequest(ctx, url); err != nil {
					cancel()
					respErr = err
					return
				}
			}
		}()
	}

	go func() {
		for i := range urls {
			ch <- urls[i]
		}
		close(ch)
	}()

	wg.Wait()
	if respErr != nil {
		return nil, respErr
	}

	return fillByOkCodes(urls), nil
}

func (c *Client) sendRequest(ctx context.Context, url string) error {
	if err := c.validateURL(url); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	httpClient := http.Client{Timeout: c.cfg.RequestTimeout}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return ErrIncorrectResponseStatusCode
	}

	return nil
}

func (c *Client) validateURL(url string) error {
	if c.cfg.AllowAll {
		return nil
	}

	if len(c.cfg.AllowedURLs) > 0 && !slices.Contains[[]string, string](c.cfg.AllowedURLs, url) {
		return ErrURLIsNotAllowed
	}

	if slices.Contains[[]string, string](c.cfg.BlockedURLs, url) {
		return ErrURLIsBlocked
	}

	return nil
}
