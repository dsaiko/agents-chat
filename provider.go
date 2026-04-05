package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicopt "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/ollama/ollama/api"
	"github.com/openai/openai-go/v3"
	openaiopt "github.com/openai/openai-go/v3/option"
)

// Provider name constants used as keys in the Providers map.
const (
	ProviderOpenAI      = "openai"
	ProviderAnthropic   = "anthropic"
	ProviderOllama      = "ollama"
	ProviderOpenRouter  = "openrouter"
)

// defaultMaxTokens is the fallback max token limit for providers that require it (e.g., Anthropic).
const defaultMaxTokens = 1024

// GenerateParams holds optional parameters for a completion request.
// Pointer fields use nil to mean "use provider default"; zero is a valid explicit value.
// MaxTokens uses 0 as "not set" since 0 tokens is not a meaningful value; providers
// that require a limit (e.g., Anthropic) fall back to defaultMaxTokens.
type GenerateParams struct {
	MaxTokens   int      // Maximum number of tokens to generate (0 = not set; Anthropic defaults to defaultMaxTokens)
	Temperature *float64 // Sampling temperature (nil = provider default, 0 = deterministic)
	TopP        *float64 // Nucleus sampling threshold (nil = provider default)
}

// Provider abstracts an LLM API for text completion.
type Provider interface {
	Generate(ctx context.Context, model string, systemPrompt string, userPrompt string, params GenerateParams) (string, error)
}

// Providers maps provider names to their implementations.
type Providers map[string]Provider

// ForModel returns the Provider and resolved model name for a given model identifier.
// Provider-specific prefixes (e.g., "ollama/") are stripped from the returned model name.
func (providers Providers) ForModel(model string) (Provider, string, error) {
	name, resolvedModel := resolveModel(model)
	p, ok := providers[name]
	if !ok {
		return nil, "", fmt.Errorf("no %s provider configured (missing API key?)", name)
	}
	return p, resolvedModel, nil
}

// initProviders creates and returns providers based on available API keys and local services.
func initProviders() Providers {
	providers := Providers{}
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		providers[ProviderOpenAI] = NewOpenAIProvider(openai.NewClient(openaiopt.WithAPIKey(key)))
	}
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		providers[ProviderAnthropic] = NewAnthropicProvider(anthropic.NewClient(anthropicopt.WithAPIKey(key)))
	}
	if client, err := api.ClientFromEnvironment(); err == nil {
		providers[ProviderOllama] = NewOllamaProvider(client)
	}
	if key := os.Getenv("OPENROUTER_API_KEY"); key != "" {
		providers[ProviderOpenRouter] = NewOpenAIProvider(openai.NewClient(
			openaiopt.WithAPIKey(key),
			openaiopt.WithBaseURL("https://openrouter.ai/api/v1/"),
		))
	}
	return providers
}

// resolveModel maps a model identifier to a provider name and the model name to pass to the API.
// Models prefixed with "ollama/" route to Ollama (prefix stripped), "openrouter/" to OpenRouter (prefix stripped),
// "claude" to Anthropic, all others to OpenAI.
func resolveModel(model string) (providerName string, resolvedModel string) {
	if m, ok := strings.CutPrefix(model, "ollama/"); ok {
		return ProviderOllama, m
	}
	if m, ok := strings.CutPrefix(model, "openrouter/"); ok {
		return ProviderOpenRouter, m
	}
	if strings.HasPrefix(model, "claude") {
		return ProviderAnthropic, model
	}
	return ProviderOpenAI, model
}
