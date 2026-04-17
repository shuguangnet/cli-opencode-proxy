package domain

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type ChatRequest struct {
	Model       string
	Messages    []Message
	Temperature *float64
	MaxTokens   *int
	Stop        []string
	Stream      bool
}

type ChatResponse struct {
	ID           string
	Model        string
	Content      string
	FinishReason string
	Usage        *Usage
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamEvent struct {
	Delta        string
	FinishReason string
	Done         bool
}
