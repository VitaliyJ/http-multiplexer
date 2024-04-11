package app

import (
	"context"
	"http-multiplexer/internal/server/http"
	"http-multiplexer/internal/service"
	httpClient "http-multiplexer/pkg/http"
)

type App struct {
	cfg *Config
}

func NewApp(cfg *Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setGlobalLogger(a.cfg.LogLevel)

	httpc := httpClient.NewClient(&httpClient.Config{AllowAll: true})
	senderService := service.NewSender(httpc)
	httpServer := http.NewServer(&http.Config{Port: a.cfg.HTTPServerPort}, senderService)

	httpServer.Run(ctx)
}
