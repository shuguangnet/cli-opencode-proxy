package app

import (
	"testing"

	"opencode-cli-proxy/internal/config"
)

func TestGatewayAddress(t *testing.T) {
	cfg := &config.Config{}
	cfg.Server.Host = "127.0.0.1"
	cfg.Server.Port = 8080
	g := NewGateway(cfg)
	if g.Address() != "127.0.0.1:8080" {
		t.Fatalf("unexpected address %s", g.Address())
	}
}
