package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"opencode-cli-proxy/internal/adapter"
	"opencode-cli-proxy/internal/config"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			status, resp := adapter.MapError(http.StatusUnauthorized, errors.New("missing bearer token"))
			c.AbortWithStatusJSON(status, resp)
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		keyCfg, ok := cfg.Keys[token]
		if !ok {
			status, resp := adapter.MapError(http.StatusUnauthorized, errors.New("invalid api key"))
			c.AbortWithStatusJSON(status, resp)
			return
		}
		c.Set("gateway_key", token)
		c.Set("account_name", keyCfg.Account)
		c.Set("allowed_models", keyCfg.AllowedModels)
		c.Next()
	}
}
