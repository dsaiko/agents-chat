package main

import (
	"context"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

// OpenAIProvider implements Provider using the OpenAI Responses API.
type OpenAIProvider struct {
	client openai.Client
}

// NewOpenAIProvider creates a new OpenAI provider with the given client.
func NewOpenAIProvider(client openai.Client) *OpenAIProvider {
	return &OpenAIProvider{client: client}
}

// Complete sends a prompt to OpenAI and returns the text response.
func (p *OpenAIProvider) Complete(ctx context.Context, model string, systemPrompt string, userPrompt string) (string, error) {
	params := responses.ResponseNewParams{
		Model: model,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(userPrompt),
		},
	}
	if systemPrompt != "" {
		params.Instructions = openai.String(systemPrompt)
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(resp.OutputText()), nil
}
