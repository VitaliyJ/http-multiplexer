package http

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	cfg     *Config
	limiter *Limiter
	sender  Sender
}

func NewServer(cfg *Config, sender Sender) *Server {
	if cfg.MaxRequestsCount == 0 {
		cfg.MaxRequestsCount = defaultMaxRequestsCount
	}

	cfg.Port = ":" + strings.TrimLeft(cfg.Port, ":")
	return &Server{
		cfg:     cfg,
		limiter: NewLimiter(cfg.MaxRequestsCount),
		sender:  sender,
	}
}

func (s *Server) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:              s.cfg.Port,
		Handler:           s.initRoutes(),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       time.Second,
		WriteTimeout:      time.Second,
	}

	httpErrCh := make(chan error, 1)
	go func() {
		httpErrCh <- srv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-httpErrCh:
		slog.Error("http server error:", err)
	case <-sigCh:
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("http server shutdown error:", err)
		}
	}
}

func (s *Server) initRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /send", s.sendHandler)

	return s.limiter.Middleware(mux)
}

func (s *Server) sendHandler(w http.ResponseWriter, r *http.Request) {
	req, err := jsonDecode[*SendRequest](r)
	if err != nil {
		jsonEncode(w, http.StatusBadRequest, &ErrorResponse{Message: "incorrect urls"})
		return
	}
	if err := req.Validate(); err != nil {
		jsonEncode(w, http.StatusBadRequest, &ErrorResponse{Message: err.Error()})
		return
	}

	res, err := s.sender.SendRequestsToURLs(r.Context(), req.URLs)
	if err != nil {
		jsonEncode(w, http.StatusBadGateway, &ErrorResponse{Message: err.Error()})
		return
	}

	jsonEncode(w, http.StatusOK, SendResponse{URLsCodes: res})
}
