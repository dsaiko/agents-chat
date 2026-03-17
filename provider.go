package main

import (
	"context"
	"fmt"
	"strings"
)

const (
	ProviderOpenAI    = "openai"
	ProviderAnthropic = "anthropic"
)

// Provider abstracts an LLM API for text completion.
type Provider interface {
	Complete(ctx context.Context, model string, systemPrompt string, userPrompt string) (string, error)
}

// providers maps provider name to its Provider implementation.
var providers = map[string]Provider{}

// providerForModel returns the Provider for a given model name.
// Models starting with "claude" use Anthropic, all others use OpenAI.
func providerForModel(model string) (Provider, error) {
	name := providerNameForModel(model)
	p, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("no %s provider configured (missing API key?)", name)
	}
	return p, nil
}

// providerNameForModel maps a model identifier to a provider name.
// Models prefixed with "claude" route to Anthropic; all others route to OpenAI.
func providerNameForModel(model string) string {
	if strings.HasPrefix(model, "claude") {
		return ProviderAnthropic
	}
	return ProviderOpenAI
}
