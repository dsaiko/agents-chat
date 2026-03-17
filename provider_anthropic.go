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

// Complete sends a prompt to Anthropic and returns the text response.
func (p *AnthropicProvider) Complete(ctx context.Context, model string, systemPrompt string, userPrompt string) (string, error) {
	params := anthropic.MessageNewParams{
		MaxTokens: 1024,
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
