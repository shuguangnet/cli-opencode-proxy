package openai

import (
	"fmt"
	"net/http"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"opencode-cli-proxy/internal/adapter"
	"opencode-cli-proxy/internal/config"
	"opencode-cli-proxy/internal/domain"
	"opencode-cli-proxy/internal/upstream"
)

type Handler struct {
	cfg    *config.Config
	client *upstream.Client
}

func NewHandler(cfg *config.Config, client *upstream.Client) *Handler {
	return &Handler{cfg: cfg, client: client}
}

func (h *Handler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":      "opencode-cli-proxy",
		"status":    "ok",
		"endpoints": []string{"/health", "/v1", "/v1/models", "/v1/chat/completions", "/v1/completions"},
	})
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) V1Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"object":    "info",
		"message":   "Use /v1/models or /v1/chat/completions.",
		"endpoints": []string{"/v1/models", "/v1/chat/completions", "/v1/completions"},
	})
}

func (h *Handler) ListModels(c *gin.Context) {
	items, err := h.listModels(c)
	if err != nil {
		h.writeError(c, http.StatusBadGateway, err)
		return
	}
	c.JSON(http.StatusOK, domain.OpenAIModelsResponse{Object: "list", Data: items})
}

func (h *Handler) listModels(c *gin.Context) ([]domain.OpenAIModelInfo, error) {
	cmdArgs := []string{"models"}
	if h.cfg.Upstream.Attach != "" {
		cmdArgs = append(cmdArgs, "--attach", h.cfg.Upstream.Attach)
	}
	cmd := exec.CommandContext(c.Request.Context(), h.cfg.Upstream.Binary, cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	items := make([]domain.OpenAIModelInfo, 0)
	seen := map[string]struct{}{}
	for _, line := range strings.Split(string(output), "\n") {
		model := strings.TrimSpace(line)
		if model == "" {
			continue
		}
		if _, ok := seen[model]; ok {
			continue
		}
		seen[model] = struct{}{}
		items = append(items, domain.OpenAIModelInfo{
			ID:      model,
			Object:  "model",
			Created: now,
			OwnedBy: "opencode",
		})
	}
	return items, nil
}

func (h *Handler) ChatCompletions(c *gin.Context) {
	var req domain.OpenAIChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.writeError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.authorizeModel(c, req.Model); err != nil {
		h.writeError(c, http.StatusForbidden, err)
		return
	}

	mapped, err := adapter.MapChatRequest(req, h.cfg)
	if err != nil {
		h.writeError(c, http.StatusBadRequest, err)
		return
	}

	if req.Stream {
		h.streamChat(c, req.Model, mapped)
		return
	}

	accountName := c.GetString("account_name")
	resp, status, err := h.client.Chat(c.Request.Context(), accountName, mapped)
	if err != nil {
		h.writeError(c, status, err)
		return
	}

	c.JSON(http.StatusOK, domain.OpenAIChatCompletionResponse{
		ID:      resp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []domain.OpenAIChoice{{
			Index: 0,
			Message: &domain.OpenAIChatMessage{
				Role:    "assistant",
				Content: resp.Content,
			},
			FinishReason: resp.FinishReason,
		}},
		Usage: resp.Usage,
	})
}

func (h *Handler) Completions(c *gin.Context) {
	var req domain.OpenAICompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.writeError(c, http.StatusBadRequest, err)
		return
	}
	if err := h.authorizeModel(c, req.Model); err != nil {
		h.writeError(c, http.StatusForbidden, err)
		return
	}

	mapped, err := adapter.MapCompletionToChat(req, h.cfg)
	if err != nil {
		h.writeError(c, http.StatusBadRequest, err)
		return
	}

	if req.Stream {
		h.streamCompletion(c, req.Model, mapped)
		return
	}

	accountName := c.GetString("account_name")
	resp, status, err := h.client.Chat(c.Request.Context(), accountName, mapped)
	if err != nil {
		h.writeError(c, status, err)
		return
	}

	c.JSON(http.StatusOK, domain.OpenAIChatCompletionResponse{
		ID:      resp.ID,
		Object:  "text_completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []domain.OpenAIChoice{{
			Index:        0,
			Text:         resp.Content,
			FinishReason: resp.FinishReason,
		}},
		Usage: resp.Usage,
	})
}

func (h *Handler) streamChat(c *gin.Context, clientModel string, req domain.ChatRequest) {
	accountName := c.GetString("account_name")
	events, errCh, status, err := h.client.ChatStream(c.Request.Context(), accountName, req)
	if err != nil {
		h.writeError(c, status, err)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Status(http.StatusOK)
	c.Writer.Flush()

	first := true
	for event := range events {
		event = adapter.NormalizeStreamEvent(event)
		if event.Done {
			_, _ = c.Writer.Write(adapter.BuildDoneChunk())
			c.Writer.Flush()
			if err, ok := <-errCh; ok && err != nil {
				return
			}
			return
		}
		chunk, err := adapter.BuildStreamChunk(fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()), clientModel, 0, event.Delta, first, event.FinishReason)
		if err != nil {
			return
		}
		first = false
		_, _ = c.Writer.Write(chunk)
		c.Writer.Flush()
	}

	if err, ok := <-errCh; ok && err != nil {
		return
	}
}

func (h *Handler) streamCompletion(c *gin.Context, clientModel string, req domain.ChatRequest) {
	h.streamChat(c, clientModel, req)
}

func (h *Handler) authorizeModel(c *gin.Context, model string) error {
	allowed, _ := c.Get("allowed_models")
	allowedModels, _ := allowed.([]string)
	if len(allowedModels) == 0 {
		return nil
	}
	if !slices.Contains(allowedModels, model) {
		return fmt.Errorf("model %s is not allowed for this api key", model)
	}
	return nil
}

func (h *Handler) writeError(c *gin.Context, status int, err error) {
	status, resp := adapter.MapError(status, err)
	c.JSON(status, resp)
}
