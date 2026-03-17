package main

import (
	"context"
	"fmt"
	"testing"
)

type mockProvider struct {
	response string
	err      error
}

func (m *mockProvider) Complete(_ context.Context, _ string, _ string, _ string) (string, error) {
	return m.response, m.err
}

func TestProviderNameForModel(t *testing.T) {
	tests := []struct {
		model string
		want  string
	}{
		{"gpt-4o", ProviderOpenAI},
		{"gpt-5-mini", ProviderOpenAI},
		{"claude-sonnet-4-5-20250514", ProviderAnthropic},
		{"claude-haiku-4-5", ProviderAnthropic},
		{"some-other-model", ProviderOpenAI},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			got := providerNameForModel(tt.model)
			if got != tt.want {
				t.Errorf("providerNameForModel(%q) = %q, want %q", tt.model, got, tt.want)
			}
		})
	}
}

func TestProviderForModel(t *testing.T) {
	// Save and restore global state
	orig := providers
	defer func() { providers = orig }()

	mock := &mockProvider{response: "ok"}
	providers = map[string]Provider{
		ProviderOpenAI: mock,
	}

	// OpenAI model should resolve
	p, err := providerForModel("gpt-4o")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != mock {
		t.Error("expected mock provider")
	}

	// Claude model should fail (no anthropic provider registered)
	_, err = providerForModel("claude-haiku-4-5")
	if err == nil {
		t.Fatal("expected error for missing anthropic provider")
	}
}

func TestRunAgentWithMock(t *testing.T) {
	orig := providers
	defer func() { providers = orig }()

	providers = map[string]Provider{
		ProviderOpenAI: &mockProvider{response: "Hello from mock"},
	}

	lang := languages["en"]
	agent := Agent{Name: "Test", Model: "gpt-4o", Instructions: "Be helpful."}
	history := []string{"Moderator: Test question"}

	reply, err := runAgent(context.Background(), lang, agent, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "Hello from mock" {
		t.Errorf("reply = %q, want %q", reply, "Hello from mock")
	}
}

func TestRunAgentEmptyReply(t *testing.T) {
	orig := providers
	defer func() { providers = orig }()

	providers = map[string]Provider{
		ProviderOpenAI: &mockProvider{response: ""},
	}

	lang := languages["en"]
	agent := Agent{Name: "Test", Model: "gpt-4o", Instructions: "Be helpful."}

	reply, err := runAgent(context.Background(), lang, agent, []string{"Moderator: Hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != lang.EmptyReply {
		t.Errorf("reply = %q, want %q", reply, lang.EmptyReply)
	}
}

func TestRunAgentProviderError(t *testing.T) {
	orig := providers
	defer func() { providers = orig }()

	providers = map[string]Provider{
		ProviderOpenAI: &mockProvider{err: fmt.Errorf("API error")},
	}

	lang := languages["en"]
	agent := Agent{Name: "Test", Model: "gpt-4o", Instructions: "Be helpful."}

	_, err := runAgent(context.Background(), lang, agent, []string{"Moderator: Hi"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBuildPrompt(t *testing.T) {
	lang := languages["en"]
	history := []string{"Moderator: Topic", "Agent A: Reply 1"}

	prompt := buildPrompt(lang, history)

	if !contains(prompt, lang.ConversationPre) {
		t.Error("prompt missing ConversationPre")
	}
	if !contains(prompt, lang.ConversationPost) {
		t.Error("prompt missing ConversationPost")
	}
	if !contains(prompt, "Moderator: Topic") {
		t.Error("prompt missing history entry")
	}
	if !contains(prompt, "Agent A: Reply 1") {
		t.Error("prompt missing history entry")
	}
}

func TestBuildPromptTruncation(t *testing.T) {
	lang := languages["en"]

	// Create 10 history entries — only last 8 should appear
	var history []string
	for i := 0; i < 10; i++ {
		history = append(history, fmt.Sprintf("Entry %d", i))
	}

	prompt := buildPrompt(lang, history)

	if contains(prompt, "Entry 0") {
		t.Error("prompt should not contain Entry 0 (truncated)")
	}
	if contains(prompt, "Entry 1") {
		t.Error("prompt should not contain Entry 1 (truncated)")
	}
	if !contains(prompt, "Entry 2") {
		t.Error("prompt should contain Entry 2")
	}
	if !contains(prompt, "Entry 9") {
		t.Error("prompt should contain Entry 9")
	}
}

func TestDetectLanguageWithMock(t *testing.T) {
	orig := providers
	defer func() { providers = orig }()

	providers = map[string]Provider{
		ProviderOpenAI: &mockProvider{response: "cs"},
	}

	lang := detectLanguage(context.Background(), "gpt-4o", "Nějaký český text")
	if lang.Language != languages["cs"].Language {
		t.Errorf("got %q, want Czech", lang.Language)
	}
}

func TestDetectLanguageFallback(t *testing.T) {
	orig := providers
	defer func() { providers = orig }()

	// Unknown language code → fallback to English
	providers = map[string]Provider{
		ProviderOpenAI: &mockProvider{response: "xx"},
	}

	lang := detectLanguage(context.Background(), "gpt-4o", "Some text")
	if lang.Language != languages["en"].Language {
		t.Errorf("got %q, want English fallback", lang.Language)
	}
}

func TestDetectLanguageProviderError(t *testing.T) {
	orig := providers
	defer func() { providers = orig }()

	providers = map[string]Provider{
		ProviderOpenAI: &mockProvider{err: fmt.Errorf("fail")},
	}

	lang := detectLanguage(context.Background(), "gpt-4o", "Some text")
	if lang.Language != languages["en"].Language {
		t.Errorf("got %q, want English fallback on error", lang.Language)
	}
}

func TestDetectLanguageNoProvider(t *testing.T) {
	orig := providers
	defer func() { providers = orig }()

	providers = map[string]Provider{}

	lang := detectLanguage(context.Background(), "gpt-4o", "Some text")
	if lang.Language != languages["en"].Language {
		t.Errorf("got %q, want English fallback when no provider", lang.Language)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
