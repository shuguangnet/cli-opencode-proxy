package server

import (
	"github.com/gin-gonic/gin"

	"opencode-cli-proxy/internal/config"
	"opencode-cli-proxy/internal/openai"
	"opencode-cli-proxy/internal/server/middleware"
)

func NewRouter(cfg *config.Config, handler *openai.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	r.GET("/", handler.Root)
	r.GET("/health", handler.Health)
	r.GET("/v1", handler.V1Root)

	v1 := r.Group("/v1")
	v1.Use(middleware.Auth(cfg))
	v1.GET("/models", handler.ListModels)
	v1.POST("/chat/completions", handler.ChatCompletions)
	v1.POST("/completions", handler.Completions)

	return r
}
