package upstream

import (
	"fmt"
	"net/http"
	"strings"

	"opencode-cli-proxy/internal/config"
)

type AuthProvider interface {
	Apply(req *http.Request) error
}

type bearerProvider struct {
	token string
}

func (p bearerProvider) Apply(req *http.Request) error {
	if strings.TrimSpace(p.token) == "" {
		return fmt.Errorf("empty bearer token")
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	return nil
}

func NewAuthProvider(account config.AccountConfig) (AuthProvider, error) {
	switch account.AuthMode {
	case "", "bearer_token":
		return bearerProvider{token: account.Token}, nil
	case "iam_token":
		return bearerProvider{token: account.Token}, nil
	case "local":
		return bearerProvider{token: "local-opencode"}, nil
	default:
		return nil, fmt.Errorf("unsupported auth mode: %s", account.AuthMode)
	}
}
