package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

type mockProvider struct {
	response string
	err      error
}

func (m *mockProvider) Generate(_ context.Context, _ string, _ string, _ string, _ GenerateParams) (string, error) {
	return m.response, m.err
}

func TestResolveModel(t *testing.T) {
	tests := []struct {
		model        string
		wantProvider string
		wantModel    string
	}{
		{"gpt-4o", ProviderOpenAI, "gpt-4o"},
		{"gpt-5-mini", ProviderOpenAI, "gpt-5-mini"},
		{"claude-sonnet-4-5-20250514", ProviderAnthropic, "claude-sonnet-4-5-20250514"},
		{"claude-haiku-4-5", ProviderAnthropic, "claude-haiku-4-5"},
		{"ollama/qwen3:8b", ProviderOllama, "qwen3:8b"},
		{"ollama/llama3", ProviderOllama, "llama3"},
		{"openrouter/qwen/qwen3.6-plus:free", ProviderOpenRouter, "qwen/qwen3.6-plus:free"},
		{"openrouter/google/gemma-2-9b-it", ProviderOpenRouter, "google/gemma-2-9b-it"},
		{"some-other-model", ProviderOpenAI, "some-other-model"},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			gotProvider, gotModel := resolveModel(tt.model)
			if gotProvider != tt.wantProvider {
				t.Errorf("resolveModel(%q) provider = %q, want %q", tt.model, gotProvider, tt.wantProvider)
			}
			if gotModel != tt.wantModel {
				t.Errorf("resolveModel(%q) model = %q, want %q", tt.model, gotModel, tt.wantModel)
			}
		})
	}
}

func TestForModel(t *testing.T) {
	mock := &mockProvider{response: "ok"}
	providers := Providers{
		ProviderOpenAI: mock,
	}

	// OpenAI model should resolve
	p, _, err := providers.ForModel("gpt-4o")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p != mock {
		t.Error("expected mock provider")
	}

	// Claude model should fail (no anthropic provider registered)
	_, _, err = providers.ForModel("claude-haiku-4-5")
	if err == nil {
		t.Fatal("expected error for missing anthropic provider")
	}

	// Ollama model should fail (no ollama provider registered)
	_, _, err = providers.ForModel("ollama/qwen3:8b")
	if err == nil {
		t.Fatal("expected error for missing ollama provider")
	}

	// OpenRouter model should fail (no openrouter provider registered)
	_, _, err = providers.ForModel("openrouter/google/gemma-2-9b-it")
	if err == nil {
		t.Fatal("expected error for missing openrouter provider")
	}
}

func TestRunAgentWithMock(t *testing.T) {
	providers := Providers{
		ProviderOpenAI: &mockProvider{response: "Hello from mock"},
	}

	lang := defaultLanguage
	agent := Agent{Name: "Test", Model: "gpt-4o", Instructions: "Be helpful."}
	history := []string{"Moderator: Test question"}

	reply, err := runAgent(context.Background(), providers, lang, agent, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "Hello from mock" {
		t.Errorf("reply = %q, want %q", reply, "Hello from mock")
	}
}

func TestRunAgentWithOllamaModel(t *testing.T) {
	providers := Providers{
		ProviderOllama: &mockProvider{response: "Hello from Ollama"},
	}

	lang := defaultLanguage
	agent := Agent{Name: "Test", Model: "ollama/qwen3:8b", Instructions: "Be helpful."}
	history := []string{"Moderator: Test question"}

	reply, err := runAgent(context.Background(), providers, lang, agent, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "Hello from Ollama" {
		t.Errorf("reply = %q, want %q", reply, "Hello from Ollama")
	}
}

func TestRunAgentWithOpenRouterModel(t *testing.T) {
	providers := Providers{
		ProviderOpenRouter: &mockProvider{response: "Hello from OpenRouter"},
	}

	lang := defaultLanguage
	agent := Agent{Name: "Test", Model: "openrouter/google/gemma-2-9b-it", Instructions: "Be helpful."}
	history := []string{"Moderator: Test question"}

	reply, err := runAgent(context.Background(), providers, lang, agent, history)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "Hello from OpenRouter" {
		t.Errorf("reply = %q, want %q", reply, "Hello from OpenRouter")
	}
}

func TestRunAgentEmptyReply(t *testing.T) {
	providers := Providers{
		ProviderOpenAI: &mockProvider{response: ""},
	}

	lang := defaultLanguage
	agent := Agent{Name: "Test", Model: "gpt-4o", Instructions: "Be helpful."}

	reply, err := runAgent(context.Background(), providers, lang, agent, []string{"Moderator: Hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != lang.EmptyReply {
		t.Errorf("reply = %q, want %q", reply, lang.EmptyReply)
	}
}

func TestRunAgentProviderError(t *testing.T) {
	providers := Providers{
		ProviderOpenAI: &mockProvider{err: fmt.Errorf("API error")},
	}

	lang := defaultLanguage
	agent := Agent{Name: "Test", Model: "gpt-4o", Instructions: "Be helpful."}

	_, err := runAgent(context.Background(), providers, lang, agent, []string{"Moderator: Hi"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBuildPrompt(t *testing.T) {
	lang := defaultLanguage
	history := []string{"Moderator: Topic", "Agent A: Reply 1"}

	prompt := buildPrompt(lang, history)

	if !strings.Contains(prompt, lang.ConversationPre) {
		t.Error("prompt missing ConversationPre")
	}
	if !strings.Contains(prompt, lang.ConversationPost) {
		t.Error("prompt missing ConversationPost")
	}
	if !strings.Contains(prompt, "Moderator: Topic") {
		t.Error("prompt missing history entry")
	}
	if !strings.Contains(prompt, "Agent A: Reply 1") {
		t.Error("prompt missing history entry")
	}
}

func TestBuildPromptTruncation(t *testing.T) {
	lang := defaultLanguage

	// Create 10 history entries — only last 8 should appear
	var history []string
	for i := range 10 {
		history = append(history, fmt.Sprintf("Entry %d", i))
	}

	prompt := buildPrompt(lang, history)

	if strings.Contains(prompt, "Entry 0") {
		t.Error("prompt should not contain Entry 0 (truncated)")
	}
	if strings.Contains(prompt, "Entry 1") {
		t.Error("prompt should not contain Entry 1 (truncated)")
	}
	if !strings.Contains(prompt, "Entry 2") {
		t.Error("prompt should contain Entry 2")
	}
	if !strings.Contains(prompt, "Entry 9") {
		t.Error("prompt should contain Entry 9")
	}
}

func TestRoundFormat(t *testing.T) {
	lang := defaultLanguage
	got := lang.Round(2, 5)
	want := "Round 2/5"
	if got != want {
		t.Errorf("Round(2, 5) = %q, want %q", got, want)
	}
}

func TestDetectLanguageWithMock(t *testing.T) {
	czJSON := `{"moderator":"Moderátor","round_format":"Kolo %d/%d","empty_reply":"(prázdná odpověď)","conversation_pre":"Dosavadní konverzace:","conversation_post":"Odpověz jako další účastník debaty. Nenapiš nic navíc mimo svou repliku.","detected_language":"Detekovaný jazyk: Čeština"}`
	providers := Providers{
		ProviderOpenAI: &mockProvider{response: czJSON},
	}

	lang := detectLanguage(context.Background(), providers, "gpt-4o", "Nějaký český text")
	if lang.Moderator != "Moderátor" {
		t.Errorf("Moderator = %q, want %q", lang.Moderator, "Moderátor")
	}
	if lang.Language != "Detekovaný jazyk: Čeština" {
		t.Errorf("Language = %q, want Czech", lang.Language)
	}
	if got := lang.Round(1, 3); got != "Kolo 1/3" {
		t.Errorf("Round(1,3) = %q, want %q", got, "Kolo 1/3")
	}
}

func TestDetectLanguageFallback(t *testing.T) {
	// Invalid JSON → fallback to English
	providers := Providers{
		ProviderOpenAI: &mockProvider{response: "not json at all"},
	}

	lang := detectLanguage(context.Background(), providers, "gpt-4o", "Some text")
	if lang.Language != defaultLanguage.Language {
		t.Errorf("got %q, want English fallback", lang.Language)
	}
}

func TestDetectLanguageProviderError(t *testing.T) {
	providers := Providers{
		ProviderOpenAI: &mockProvider{err: fmt.Errorf("fail")},
	}

	lang := detectLanguage(context.Background(), providers, "gpt-4o", "Some text")
	if lang.Language != defaultLanguage.Language {
		t.Errorf("got %q, want English fallback on error", lang.Language)
	}
}

func TestDetectLanguageNoProvider(t *testing.T) {
	providers := Providers{}

	lang := detectLanguage(context.Background(), providers, "gpt-4o", "Some text")
	if lang.Language != defaultLanguage.Language {
		t.Errorf("got %q, want English fallback when no provider", lang.Language)
	}
}

func TestParseLanguageJSON(t *testing.T) {
	// Plain JSON
	json := `{"moderator":"Mod","round_format":"R %d/%d","empty_reply":"(e)","conversation_pre":"Pre:","conversation_post":"Post.","detected_language":"Lang: Test"}`
	lang := parseLanguageJSON(json)
	if lang.Moderator != "Mod" {
		t.Errorf("Moderator = %q, want %q", lang.Moderator, "Mod")
	}

	// JSON wrapped in markdown code fences
	fenced := "```json\n" + json + "\n```"
	lang = parseLanguageJSON(fenced)
	if lang.Moderator != "Mod" {
		t.Errorf("fenced: Moderator = %q, want %q", lang.Moderator, "Mod")
	}

	// Missing round_format %d → should use default
	noFmt := `{"moderator":"M","round_format":"Runde","empty_reply":"(e)","conversation_pre":"P:","conversation_post":"P.","detected_language":"L"}`
	lang = parseLanguageJSON(noFmt)
	if lang.RoundFormat != defaultLanguage.RoundFormat {
		t.Errorf("RoundFormat = %q, want default %q", lang.RoundFormat, defaultLanguage.RoundFormat)
	}

	// Invalid JSON → default
	lang = parseLanguageJSON("garbage")
	if lang != defaultLanguage {
		t.Errorf("expected defaultLanguage for invalid JSON")
	}
}
