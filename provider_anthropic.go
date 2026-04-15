package main

import (
	"context"
	"fmt"
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

// HealthCheck validates credentials and base URL against a lightweight models API call.
func (p *AnthropicProvider) HealthCheck(ctx context.Context) error {
	_, err := p.client.Models.List(ctx, anthropic.ModelListParams{Limit: anthropic.Int(1)})
	if err != nil {
		return fmt.Errorf("anthropic health check failed: %w", err)
	}
	return nil
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
	if cp.Temperature != nil {
		params.Temperature = anthropic.Float(*cp.Temperature)
	}
	if cp.TopP != nil {
		params.TopP = anthropic.Float(*cp.TopP)
	}

	msg, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return "", err
	}

	switch msg.StopReason {
	case "", anthropic.StopReasonEndTurn, anthropic.StopReasonStopSequence:
		// OK.
	case anthropic.StopReasonMaxTokens:
		return "", fmt.Errorf("response incomplete: stopped at max_tokens")
	case anthropic.StopReasonRefusal:
		return "", fmt.Errorf("response refused")
	default:
		return "", fmt.Errorf("unsupported stop reason: %s", msg.StopReason)
	}

	var b strings.Builder
	for _, block := range msg.Content {
		if block.Type == "text" {
			b.WriteString(block.AsText().Text)
		}
	}

	text := strings.TrimSpace(b.String())
	if text == "" && len(msg.Content) > 0 {
		return "", fmt.Errorf("response contained no text output")
	}

	return text, nil
}
