package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Language holds localized UI strings for a specific language.
type Language struct {
	Moderator        string `json:"moderator"`
	RoundFormat      string `json:"round_format"`
	EmptyReply       string `json:"empty_reply"`
	ConversationPre  string `json:"conversation_pre"`
	ConversationPost string `json:"conversation_post"`
	Language         string `json:"detected_language"`
}

// Round formats the round header using the localized format string.
func (l Language) Round(i, total int) string {
	return fmt.Sprintf(l.RoundFormat, i, total)
}

// defaultLanguage provides English defaults used as fallback.
var defaultLanguage = Language{
	Moderator:        "Moderator",
	RoundFormat:      "Round %d/%d",
	EmptyReply:       "(empty reply)",
	ConversationPre:  "Conversation so far:",
	ConversationPost: "Reply as the next participant of the debate. Write nothing beyond your reply.",
	Language:         "Detected language: English",
}

// detectLanguage uses an LLM to detect the language of text and translate UI strings.
// Falls back to English defaults on any failure.
func detectLanguage(ctx context.Context, model string, text string) Language {
	p, err := providerForModel(model)
	if err != nil {
		return defaultLanguage
	}

	originals, _ := json.MarshalIndent(defaultLanguage, "", "  ")

	prompt := fmt.Sprintf(`Detect the language of the text below and translate these UI strings into that language.
Reply with ONLY a valid JSON object — no markdown, no code fences, no commentary.

English originals:
%s

IMPORTANT:
- Keep the two %%d/%%d placeholders exactly as-is in "round_format"
- In "detected_language", translate the phrase and use the actual language name

Text:
%s`, originals, text)

	result, err := p.Complete(ctx, model, "", prompt)
	if err != nil {
		return defaultLanguage
	}

	return parseLanguageJSON(result)
}

// parseLanguageJSON extracts a Language from an LLM response, stripping markdown fences if present.
func parseLanguageJSON(raw string) Language {
	s := strings.TrimSpace(raw)

	// Strip markdown code fences if present
	if strings.HasPrefix(s, "```") {
		if idx := strings.Index(s[3:], "\n"); idx >= 0 {
			s = s[3+idx+1:]
		}
		if idx := strings.LastIndex(s, "```"); idx >= 0 {
			s = s[:idx]
		}
		s = strings.TrimSpace(s)
	}

	var lang Language
	if err := json.Unmarshal([]byte(s), &lang); err != nil {
		return defaultLanguage
	}

	// Validate that RoundFormat contains %d placeholders
	if lang.RoundFormat == "" || !strings.Contains(lang.RoundFormat, "%d") {
		lang.RoundFormat = defaultLanguage.RoundFormat
	}

	return lang
}
