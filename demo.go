package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Demo holds the configuration for a debate session loaded from a directory of markdown files.
type Demo struct {
	Question string  // The debate topic, loaded from Question.md body
	Rounds   int     // Number of debate rounds, from Question.md frontmatter (default 5)
	Agents   []Agent // Participating agents, sorted by name
}

// Agent represents a single debate participant with its LLM configuration.
type Agent struct {
	Name         string // Display name from frontmatter
	Model        string // LLM model identifier (determines which provider to use)
	MaxTokens    int    // Maximum response tokens (0 = provider default)
	Instructions string // System prompt / personality from the markdown body
}

// Load reads a demo directory containing Question.md and agent .md files.
// Question.md may have optional frontmatter with "rounds" (defaults to 5).
// All other .md files are parsed as agents (must have "name" and "model" in frontmatter).
func (d *Demo) Load(dir string) error {
	questionData, err := os.ReadFile(filepath.Join(dir, "Question.md"))
	if err != nil {
		return fmt.Errorf("reading Question.md: %w", err)
	}

	fm, body, err := parseFrontmatter(string(questionData))
	if err != nil {
		// No frontmatter — treat entire file as question text, default rounds
		d.Question = strings.TrimSpace(string(questionData))
		d.Rounds = 5
	} else {
		d.Question = body
		d.Rounds = 5
		if v := fm["rounds"]; v != "" {
			d.Rounds, err = strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("invalid 'rounds' in Question.md: %w", err)
			}
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory: %w", err)
	}

	d.Agents = nil
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".md") || strings.EqualFold(name, "Question.md") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("reading %s: %w", name, err)
		}

		agent, err := parseAgentFile(string(data))
		if err != nil {
			return fmt.Errorf("parsing %s: %w", name, err)
		}
		d.Agents = append(d.Agents, agent)
	}

	sort.Slice(d.Agents, func(i, j int) bool {
		return d.Agents[i].Name < d.Agents[j].Name
	})

	return nil
}

// parseAgentFile parses a markdown file with frontmatter containing "name" and "model" fields.
// The markdown body becomes the agent's system prompt / instructions.
func parseAgentFile(content string) (Agent, error) {
	fm, body, err := parseFrontmatter(content)
	if err != nil {
		return Agent{}, err
	}
	if fm["name"] == "" {
		return Agent{}, fmt.Errorf("missing 'name' in frontmatter")
	}
	if fm["model"] == "" {
		return Agent{}, fmt.Errorf("missing 'model' in frontmatter")
	}
	agent := Agent{Name: fm["name"], Model: fm["model"], Instructions: body}
	if v := fm["max_tokens"]; v != "" {
		agent.MaxTokens, err = strconv.Atoi(v)
		if err != nil {
			return Agent{}, fmt.Errorf("invalid 'max_tokens' in frontmatter: %w", err)
		}
	}
	return agent, nil
}

// parseFrontmatter parses a markdown file with YAML-like frontmatter (key: value pairs)
// and returns the frontmatter fields as a map and the body text.
func parseFrontmatter(content string) (map[string]string, string, error) {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "---") {
		return nil, "", fmt.Errorf("missing frontmatter")
	}

	content = content[3:]
	frontmatter, rest, found := strings.Cut(content, "---")
	if !found {
		return nil, "", fmt.Errorf("missing closing frontmatter delimiter")
	}

	fields := make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(frontmatter), "\n") {
		if k, v, ok := strings.Cut(strings.TrimSpace(line), ":"); ok {
			fields[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}

	body := strings.TrimSpace(rest)
	return fields, body, nil
}
