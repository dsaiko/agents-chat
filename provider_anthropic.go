package main

import (
	"context"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

// AnthropicProvider implements Provider using the Anthropic Messages API.
type AnthropicProvider struct {
	client anthropic.Client
}

// NewAnthropicProvider creates a new Anthropic provider with the given client.
func NewAnthropicProvider(client anthropic.Client) *AnthropicProvider {
	return &AnthropicProvider{client: client}
}

// Generate sends a prompt to Anthropic and returns the text response.
func (p *AnthropicProvider) Generate(ctx context.Context, model string, systemPrompt string, userPrompt string, cp GenerateParams) (string, error) {
	maxTokens := int64(cp.MaxTokens)
	if maxTokens <= 0 {
		maxTokens = int64(defaultMaxTokens)
	}
	params := anthropic.MessageNewParams{
		MaxTokens: maxTokens,
		Model:     anthropic.Model(model),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(userPrompt),
			),
		},
	}
	if systemPrompt != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: systemPrompt},
		}
	}
	if cp.Temperature > 0 {
		params.Temperature = anthropic.Float(cp.Temperature)
	}
	if cp.TopP > 0 {
		params.TopP = anthropic.Float(cp.TopP)
	}

	msg, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for _, block := range msg.Content {
		if block.Type == "text" {
			b.WriteString(block.AsText().Text)
		}
	}

	return strings.TrimSpace(b.String()), nil
}
