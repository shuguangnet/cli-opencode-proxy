package adapter

import (
	"fmt"
	"math"
	"strings"

	"opencode-cli-proxy/internal/config"
	"opencode-cli-proxy/internal/domain"
)

func MapChatRequest(req domain.OpenAIChatCompletionRequest, cfg *config.Config) (domain.ChatRequest, error) {
	upstreamModel, ok := cfg.Models[req.Model]
	if !ok {
		return domain.ChatRequest{}, fmt.Errorf("unknown model: %s", req.Model)
	}

	messages := make([]domain.Message, 0, len(req.Messages))
	for _, msg := range req.Messages {
		role := strings.TrimSpace(msg.Role)
		if role == "" {
			return domain.ChatRequest{}, fmt.Errorf("message role is required")
		}
		messages = append(messages, domain.Message{Role: role, Text: msg.Content})
	}

	mapped := domain.ChatRequest{
		Model:       upstreamModel,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stop:        req.Stop,
		Stream:      req.Stream,
	}

	if req.Temperature != nil {
		v := MapTemperature(*req.Temperature, cfg.Mapping.Temperature.TargetMin, cfg.Mapping.Temperature.TargetMax)
		mapped.Temperature = &v
	}

	return mapped, nil
}

func MapCompletionToChat(req domain.OpenAICompletionRequest, cfg *config.Config) (domain.ChatRequest, error) {
	chatReq := domain.OpenAIChatCompletionRequest{
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stop:        req.Stop,
		Stream:      req.Stream,
		Messages: []domain.OpenAIChatMessage{{
			Role:    "user",
			Content: req.Prompt,
		}},
	}
	return MapChatRequest(chatReq, cfg)
}

func MapTemperature(v, targetMin, targetMax float64) float64 {
	if v < 0 {
		v = 0
	}
	if v > 2 {
		v = 2
	}
	ratio := v / 2
	mapped := targetMin + ratio*(targetMax-targetMin)
	return math.Max(targetMin, math.Min(targetMax, mapped))
}
