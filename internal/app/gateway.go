package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"opencode-cli-proxy/internal/config"
	"opencode-cli-proxy/internal/openai"
	"opencode-cli-proxy/internal/server"
	"opencode-cli-proxy/internal/upstream"
)

type Gateway struct {
	cfg        *config.Config
	httpServer *http.Server
}

func NewGateway(cfg *config.Config) *Gateway {
	client := upstream.NewClient(cfg, nil)
	handler := openai.NewHandler(cfg, client)
	router := server.NewRouter(cfg, handler)
	return &Gateway{
		cfg:        cfg,
		httpServer: newHTTPServer(cfg, router),
	}
}

func newHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
	}
}

func (g *Gateway) Start() error {
	return g.httpServer.ListenAndServe()
}

func (g *Gateway) Stop(ctx context.Context) error {
	return g.httpServer.Shutdown(ctx)
}

func (g *Gateway) Address() string {
	return g.httpServer.Addr
}
