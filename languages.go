package main

import (
	"context"
	"fmt"
	"strings"
)

// Language holds localized UI strings for a specific language.
type Language struct {
	Moderator        string                 // Label for the moderator in conversation history
	Round            func(i, total int) string // Formats the round header (e.g. "Round 1/5")
	EmptyReply       string                 // Placeholder when an agent returns no text
	ConversationPre  string                 // Prompt prefix before conversation history
	ConversationPost string                 // Prompt suffix instructing the agent to reply
	Language         string                 // Human-readable language detection result
}

// languages maps ISO 639-1 codes to localized UI strings.
var languages = map[string]Language{
	"cs": {
		Moderator:        "Moderátor",
		Round:            func(i, total int) string { return fmt.Sprintf("Kolo %d/%d", i, total) },
		EmptyReply:       "(prázdná odpověď)",
		ConversationPre:  "Dosavadní konverzace:",
		ConversationPost: "Odpověz jako další účastník debaty. Nenapiš nic navíc mimo svou repliku.",
		Language:         "Detekovaný jazyk: Čeština",
	},
	"en": {
		Moderator:        "Moderator",
		Round:            func(i, total int) string { return fmt.Sprintf("Round %d/%d", i, total) },
		EmptyReply:       "(empty reply)",
		ConversationPre:  "Conversation so far:",
		ConversationPost: "Reply as the next participant of the debate. Write nothing beyond your reply.",
		Language:         "Detected language: English",
	},
	"de": {
		Moderator:        "Moderator",
		Round:            func(i, total int) string { return fmt.Sprintf("Runde %d/%d", i, total) },
		EmptyReply:       "(leere Antwort)",
		ConversationPre:  "Bisheriges Gespräch:",
		ConversationPost: "Antworte als nächster Teilnehmer der Debatte. Schreibe nichts außer deiner Antwort.",
		Language:         "Erkannte Sprache: Deutsch",
	},
	"fr": {
		Moderator:        "Modérateur",
		Round:            func(i, total int) string { return fmt.Sprintf("Tour %d/%d", i, total) },
		EmptyReply:       "(réponse vide)",
		ConversationPre:  "Conversation jusqu'ici :",
		ConversationPost: "Réponds en tant que prochain participant du débat. N'écris rien d'autre que ta réplique.",
		Language:         "Langue détectée : Français",
	},
	"es": {
		Moderator:        "Moderador",
		Round:            func(i, total int) string { return fmt.Sprintf("Ronda %d/%d", i, total) },
		EmptyReply:       "(respuesta vacía)",
		ConversationPre:  "Conversación hasta ahora:",
		ConversationPost: "Responde como el siguiente participante del debate. No escribas nada más que tu réplica.",
		Language:         "Idioma detectado: Español",
	},
	"pt": {
		Moderator:        "Moderador",
		Round:            func(i, total int) string { return fmt.Sprintf("Rodada %d/%d", i, total) },
		EmptyReply:       "(resposta vazia)",
		ConversationPre:  "Conversa até agora:",
		ConversationPost: "Responda como o próximo participante do debate. Não escreva nada além da sua réplica.",
		Language:         "Idioma detectado: Português",
	},
	"it": {
		Moderator:        "Moderatore",
		Round:            func(i, total int) string { return fmt.Sprintf("Turno %d/%d", i, total) },
		EmptyReply:       "(risposta vuota)",
		ConversationPre:  "Conversazione finora:",
		ConversationPost: "Rispondi come il prossimo partecipante del dibattito. Non scrivere nient'altro che la tua replica.",
		Language:         "Lingua rilevata: Italiano",
	},
}

// detectLanguage uses an LLM provider to detect the language of the given text
// and returns the matching Language. Falls back to "en" on any failure.
func detectLanguage(ctx context.Context, model string, text string) Language {
	p, err := providerForModel(model)
	if err != nil {
		return languages["en"]
	}

	result, err := p.Complete(ctx, model, "", "What language is the following text? Reply with ONLY the ISO 639-1 two-letter code, nothing else.\n\n"+text)
	if err != nil {
		return languages["en"]
	}

	code := strings.TrimSpace(strings.ToLower(result))
	if lang, ok := languages[code]; ok {
		return lang
	}

	return languages["en"]
}
