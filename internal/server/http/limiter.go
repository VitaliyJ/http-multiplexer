package http

import (
	"errors"
	"net/http"
	"sync/atomic"
)

type Limiter struct {
	max int64
	c   atomic.Int64
}

func NewLimiter(max int64) *Limiter {
	return &Limiter{
		max: max,
	}
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer l.Release()
		if err := l.Can(); err != nil {
			jsonEncode(w, http.StatusTooManyRequests, &ErrorResponse{Message: err.Error()})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (l *Limiter) Can() error {
	l.c.Add(1)
	if l.c.Load() > l.max {
		return errors.New("limit exceeded")
	}

	return nil
}

func (l *Limiter) Release() {
	l.c.Add(-1)
}
