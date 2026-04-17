package upstream

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"opencode-cli-proxy/internal/config"
	"opencode-cli-proxy/internal/domain"
)

type Client struct {
	cfg *config.Config
}

type opencodeTextEvent struct {
	Type string `json:"type"`
	Part struct {
		Text string `json:"text"`
	} `json:"part"`
}

type opencodeFinishEvent struct {
	Type string `json:"type"`
	Part struct {
		Tokens struct {
			Input     int `json:"input"`
			Output    int `json:"output"`
			Total     int `json:"total"`
			Reasoning int `json:"reasoning"`
		} `json:"tokens"`
	} `json:"part"`
}

func NewClient(cfg *config.Config, _ any) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Chat(ctx context.Context, _ string, req domain.ChatRequest) (*domain.ChatResponse, int, error) {
	cmd := c.buildCommand(ctx, req)
	output, err := cmd.Output()
	if err != nil {
		return nil, 502, commandError(err)
	}

	resp, err := parseOpencodeOutput(strings.NewReader(string(output)), false)
	if err != nil {
		return nil, 502, err
	}
	return resp, 200, nil
}

func (c *Client) ChatStream(ctx context.Context, _ string, req domain.ChatRequest) (<-chan domain.StreamEvent, <-chan error, int, error) {
	cmd := c.buildCommand(ctx, req)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, 500, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, 500, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, 502, commandError(err)
	}

	events := make(chan domain.StreamEvent)
	errCh := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errCh)
		defer stdout.Close()
		defer stderr.Close()

		if err := streamOpencodeOutput(stdout, events); err != nil {
			errCh <- err
			_ = cmd.Wait()
			return
		}
		if err := cmd.Wait(); err != nil {
			msg, _ := io.ReadAll(stderr)
			if len(strings.TrimSpace(string(msg))) > 0 {
				errCh <- fmt.Errorf(strings.TrimSpace(string(msg)))
				return
			}
			errCh <- commandError(err)
		}
	}()

	return events, errCh, 200, nil
}

func (c *Client) buildCommand(ctx context.Context, req domain.ChatRequest) *exec.Cmd {
	args := []string{"run", "--format", "json", "--model", req.Model}
	if c.cfg.Upstream.Attach != "" {
		args = append(args, "--attach", c.cfg.Upstream.Attach)
	}
	args = append(args, buildPrompt(req))
	return exec.CommandContext(ctx, c.cfg.Upstream.Binary, args...)
}

func buildPrompt(req domain.ChatRequest) string {
	var parts []string
	for _, msg := range req.Messages {
		role := strings.ToUpper(strings.TrimSpace(msg.Role))
		text := strings.TrimSpace(msg.Text)
		if text == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("[%s]\n%s", role, text))
	}
	return strings.Join(parts, "\n\n")
}

func parseOpencodeOutput(r io.Reader, emitStream bool) (*domain.ChatResponse, error) {
	scanner := bufio.NewScanner(r)
	resp := &domain.ChatResponse{}
	var text strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var textEvent opencodeTextEvent
		if err := json.Unmarshal([]byte(line), &textEvent); err != nil {
			return nil, err
		}
		if textEvent.Type == "text" {
			text.WriteString(textEvent.Part.Text)
			continue
		}

		if textEvent.Type == "step_finish" {
			var finish opencodeFinishEvent
			if err := json.Unmarshal([]byte(line), &finish); err != nil {
				return nil, err
			}
			resp.FinishReason = "stop"
			resp.Usage = &domain.Usage{
				PromptTokens:     finish.Part.Tokens.Input,
				CompletionTokens: finish.Part.Tokens.Output,
				TotalTokens:      finish.Part.Tokens.Total,
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	resp.Content = text.String()
	if resp.FinishReason == "" && !emitStream {
		resp.FinishReason = "stop"
	}
	return resp, nil
}

func streamOpencodeOutput(r io.Reader, events chan<- domain.StreamEvent) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var event opencodeTextEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return err
		}
		switch event.Type {
		case "text":
			if event.Part.Text != "" {
				events <- domain.StreamEvent{Delta: event.Part.Text}
			}
		case "step_finish":
			events <- domain.StreamEvent{FinishReason: "stop", Done: true}
			return nil
		}
	}
	return scanner.Err()
}

func commandError(err error) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		msg := strings.TrimSpace(string(exitErr.Stderr))
		if msg != "" {
			return fmt.Errorf(msg)
		}
	}
	return err
}
