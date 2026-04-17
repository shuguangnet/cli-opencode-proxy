package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"opencode-cli-proxy/internal/app"
	"opencode-cli-proxy/internal/config"
)

func main() {
	configPath := flag.String("config", "configs/config.example.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gateway := app.NewGateway(cfg)
	log.Printf("gateway listening on %s", gateway.Address())
	if err := gateway.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("run server: %v", err)
	}
}
