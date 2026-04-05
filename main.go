package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var demoDir string
	switch {
	case len(os.Args) > 1:
		demoDir = os.Args[1]
	case os.Getenv("DEMO_DIR") != "":
		demoDir = filepath.Join("demos", os.Getenv("DEMO_DIR"))
	default:
		log.Fatal("usage: agents-chat <demo-dir> or set DEMO_DIR")
	}

	var demo Demo
	if err := demo.Load(demoDir); err != nil {
		log.Fatalf("failed to load demo from %s: %v", demoDir, err)
	}

	if len(demo.Agents) < 2 {
		log.Fatal("need at least 2 agent files")
	}

	providers := initProviders()

	// Verify all agents have a working provider
	for _, agent := range demo.Agents {
		if _, _, err := providers.ForModel(agent.Model); err != nil {
			log.Fatalf("agent %s (model %s): %v", agent.Name, agent.Model, err)
		}
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer ctxCancel()

	separator := strings.Repeat("─", 60)

	fmt.Println(separator)
	fmt.Printf("  Language detection (%s)....\n", demo.Agents[0].Model)
	lang := detectLanguage(ctx, providers, demo.Agents[0].Model, demo.Question)
	fmt.Printf("  %s\n", lang.Language)
	fmt.Printf("  %s\n", demo.Question)
	fmt.Println(separator)

	history := []string{
		lang.Moderator + ": " + demo.Question,
	}

	for i := 1; i <= demo.Rounds; i++ {
		fmt.Printf("\n── %s ──\n", lang.Round(i, demo.Rounds))
		for _, agent := range demo.Agents {
			reply, err := runAgent(ctx, providers, lang, agent, history)
			if err != nil {
				log.Fatalf("%s failed in round %d: %v", agent.Name, i, err)
			}
			indented := "    " + strings.ReplaceAll(reply, "\n", "\n    ")
			fmt.Printf("\n  [%s] (%s)\n%s\n", agent.Name, agent.Model, indented)
			history = append(history, agent.Name+": "+reply)
		}
	}

	fmt.Printf("\n%s\n", separator)
}

// runAgent sends the conversation history to the agent's LLM provider and returns the reply.
func runAgent(ctx context.Context, providers Providers, lang Language, agent Agent, history []string) (string, error) {
	p, model, err := providers.ForModel(agent.Model)
	if err != nil {
		return "", err
	}

	prompt := buildPrompt(lang, history)
	text, err := p.Generate(ctx, model, strings.TrimSpace(agent.Instructions), prompt, GenerateParams{
		MaxTokens:   agent.MaxTokens,
		Temperature: agent.Temperature,
		TopP:        agent.TopP,
	})
	if err != nil {
		return "", err
	}

	if text == "" {
		return lang.EmptyReply, nil
	}
	return text, nil
}

// buildPrompt constructs the user prompt from conversation history,
// keeping only the last maxHistory entries to stay within context limits.
func buildPrompt(lang Language, history []string) string {
	const maxHistory = 8

	start := 0
	if len(history) > maxHistory {
		start = len(history) - maxHistory
	}

	var b strings.Builder
	b.WriteString(lang.ConversationPre)
	b.WriteString("\n\n")
	for _, line := range history[start:] {
		b.WriteString(line)
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(lang.ConversationPost)
	return b.String()
}
