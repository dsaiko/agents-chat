package main

import (
	"context"
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
	if cp.Temperature > 0 {
		opts["temperature"] = cp.Temperature
	}
	if cp.TopP > 0 {
		opts["top_p"] = cp.TopP
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

	return strings.TrimSpace(response.Message.Content), nil
}
