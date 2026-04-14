package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// OpenRouterProvider implements Provider using the OpenRouter API (OpenAI-compatible).
type OpenRouterProvider struct {
	client openai.Client
}

// NewOpenRouterProvider creates a new OpenRouter provider with the given client.
func NewOpenRouterProvider(client openai.Client) *OpenRouterProvider {
	return &OpenRouterProvider{client: client}
}

// HealthCheck validates credentials and base URL against a lightweight models API call.
func (p *OpenRouterProvider) HealthCheck(ctx context.Context) error {
	if _, err := p.client.Models.List(ctx); err != nil {
		return fmt.Errorf("openrouter health check failed: %w", err)
	}
	return nil
}

// Generate sends a prompt to OpenRouter and returns the text response.
func (p *OpenRouterProvider) Generate(ctx context.Context, model string, systemPrompt string, userPrompt string, cp GenerateParams) (string, error) {
	params := responses.ResponseNewParams{
		Model: model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(userPrompt),
		},
	}
	if systemPrompt != "" {
		params.Instructions = openai.String(systemPrompt)
	}
	if cp.MaxTokens > 0 {
		params.MaxOutputTokens = openai.Int(int64(cp.MaxTokens))
	}
	if cp.Temperature != nil {
		params.Temperature = openai.Float(*cp.Temperature)
	}
	if cp.TopP != nil {
		params.TopP = openai.Float(*cp.TopP)
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return "", err
	}

	switch resp.Status {
	case "", responses.ResponseStatusCompleted:
		// OpenRouter-compatible implementations may omit the status field.
	case responses.ResponseStatusFailed:
		if resp.Error.Message != "" {
			return "", fmt.Errorf("response failed: %s", resp.Error.Message)
		}
		return "", fmt.Errorf("response failed")
	case responses.ResponseStatusIncomplete:
		if resp.IncompleteDetails.Reason != "" {
			return "", fmt.Errorf("response incomplete: %s", resp.IncompleteDetails.Reason)
		}
		return "", fmt.Errorf("response incomplete")
	default:
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	text := strings.TrimSpace(resp.OutputText())
	if text == "" && len(resp.Output) > 0 {
		return "", fmt.Errorf("response contained no text output")
	}

	return text, nil
}
