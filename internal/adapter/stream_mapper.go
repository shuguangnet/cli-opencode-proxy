package adapter

import (
	"encoding/json"
	"fmt"
	"time"

	"opencode-cli-proxy/internal/domain"
)

func BuildStreamChunk(id, model string, index int, delta string, includeRole bool, finishReason string) ([]byte, error) {
	msg := &domain.OpenAIChatMessage{Content: delta}
	if includeRole {
		msg.Role = "assistant"
	}

	resp := domain.OpenAIChatCompletionResponse{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []domain.OpenAIChoice{{
			Index:        index,
			Delta:        msg,
			FinishReason: finishReason,
		}},
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("data: %s\n\n", payload)), nil
}

func BuildDoneChunk() []byte {
	return []byte("data: [DONE]\n\n")
}

func NormalizeStreamEvent(event domain.StreamEvent) domain.StreamEvent {
	return event
}
