package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

// OllamaProvider implements Provider using the Ollama chat API.
type OllamaProvider struct {
	client *api.Client
}

// NewOllamaProvider creates a new Ollama provider with the given client.
func NewOllamaProvider(client *api.Client) *OllamaProvider {
	return &OllamaProvider{client: client}
}

// HealthCheck verifies that the Ollama server is reachable before the first turn.
func (p *OllamaProvider) HealthCheck(ctx context.Context) error {
	if err := p.client.Heartbeat(ctx); err != nil {
		return fmt.Errorf("ollama health check failed: %w", err)
	}
	return nil
}

// Generate sends a prompt to Ollama and returns the text response.
func (p *OllamaProvider) Generate(ctx context.Context, model string, systemPrompt string, userPrompt string, cp GenerateParams) (string, error) {
	var messages []api.Message
	if systemPrompt != "" {
		messages = append(messages, api.Message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, api.Message{Role: "user", Content: userPrompt})

	stream := false
	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   &stream,
	}
	opts := map[string]any{}
	if cp.MaxTokens > 0 {
		opts["num_predict"] = cp.MaxTokens
	}
	if cp.Temperature != nil {
		opts["temperature"] = *cp.Temperature
	}
	if cp.TopP != nil {
		opts["top_p"] = *cp.TopP
	}
	if len(opts) > 0 {
		req.Options = opts
	}

	var response api.ChatResponse
	err := p.client.Chat(ctx, req, func(cr api.ChatResponse) error {
		response = cr
		return nil
	})
	if err != nil {
		return "", err
	}

	text := strings.TrimSpace(response.Message.Content)
	if text == "" && len(response.Message.ToolCalls) > 0 {
		return "", fmt.Errorf("response contained tool calls instead of text output")
	}

	return text, nil
}
