package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicopt "github.com/anthropics/anthropic-sdk-go/option"
	"github.com/openai/openai-go/v3"
	openaiopt "github.com/openai/openai-go/v3/option"
)

func main() {
	demoName := os.Getenv("DEMO_DIR")
	if demoName == "" {
		log.Fatal("missing DEMO_DIR")
	}
	demoDir := filepath.Join("demos", demoName)
	if len(os.Args) > 1 {
		demoDir = os.Args[1]
	}

	var demo Demo
	if err := demo.Load(demoDir); err != nil {
		log.Fatalf("failed to load demo from %s: %v", demoDir, err)
	}

	if len(demo.Agents) < 2 {
		log.Fatal("need at least 2 agent files")
	}

	initProviders()

	// Verify all agents have a working provider
	for _, agent := range demo.Agents {
		if _, err := providerForModel(agent.Model); err != nil {
			log.Fatalf("agent %s (model %s): %v", agent.Name, agent.Model, err)
		}
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer ctxCancel()

	separator := strings.Repeat("─", 60)

	fmt.Println(separator)
	fmt.Printf("  Language detection ....\n")
	lang := detectLanguage(ctx, demo.Agents[0].Model, demo.Question)
	fmt.Printf("  %s\n", lang.Language)
	fmt.Printf("  %s\n", demo.Question)
	fmt.Println(separator)

	history := []string{
		lang.Moderator + ": " + demo.Question,
	}

	for i := 1; i <= demo.Rounds; i++ {
		fmt.Printf("\n── %s ──\n", lang.Round(i, demo.Rounds))
		for _, agent := range demo.Agents {
			reply, err := runAgent(ctx, lang, agent, history)
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

// initProviders registers LLM providers based on available API keys in environment variables.
func initProviders() {
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		providers[ProviderOpenAI] = NewOpenAIProvider(openai.NewClient(openaiopt.WithAPIKey(key)))
	}
	if key := os.Getenv("ANTHROPIC_API_KEY"); key != "" {
		providers[ProviderAnthropic] = NewAnthropicProvider(anthropic.NewClient(anthropicopt.WithAPIKey(key)))
	}
}

// runAgent sends the conversation history to the agent's LLM provider and returns the reply.
func runAgent(ctx context.Context, lang Language, agent Agent, history []string) (string, error) {
	p, err := providerForModel(agent.Model)
	if err != nil {
		return "", err
	}

	prompt := buildPrompt(lang, history)
	text, err := p.Complete(ctx, agent.Model, strings.TrimSpace(agent.Instructions), prompt)
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
