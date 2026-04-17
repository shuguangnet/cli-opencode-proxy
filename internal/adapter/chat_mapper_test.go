package adapter

import (
	"testing"

	"opencode-cli-proxy/internal/config"
	"opencode-cli-proxy/internal/domain"
)

func TestMapTemperature(t *testing.T) {
	got := MapTemperature(1, 0, 1)
	if got != 0.5 {
		t.Fatalf("expected 0.5, got %v", got)
	}
}

func TestMapChatRequest(t *testing.T) {
	cfg := &config.Config{Models: map[string]string{"opencode-go-latest": "upstream-model"}}
	req := domain.OpenAIChatCompletionRequest{
		Model: "opencode-go-latest",
		Messages: []domain.OpenAIChatMessage{{Role: "user", Content: "hello"}},
	}
	mapped, err := MapChatRequest(req, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if mapped.Model != "upstream-model" {
		t.Fatalf("unexpected model %s", mapped.Model)
	}
	if len(mapped.Messages) != 1 || mapped.Messages[0].Text != "hello" {
		t.Fatalf("unexpected messages %#v", mapped.Messages)
	}
}
