package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Demo holds the configuration for a debate session loaded from a directory of YAML files.
type Demo struct {
	Question string  // The debate topic, loaded from question.yaml
	Rounds   int     // Number of debate rounds, from question.yaml (default 5)
	Agents   []Agent // Participating agents, sorted by name
}

// Agent represents a single debate participant with its LLM configuration.
type Agent struct {
	Name         string   `yaml:"name"`
	Model        string   `yaml:"model"`
	MaxTokens    int      `yaml:"max_tokens"`
	Temperature  *float64 `yaml:"temperature"`
	TopP         *float64 `yaml:"top_p"`
	Instructions string   `yaml:"instructions"`
}

// questionFile represents the YAML structure of question.yaml.
type questionFile struct {
	Rounds   int    `yaml:"rounds"`
	Question string `yaml:"question"`
}

// Load reads a demo directory containing question.yaml and agent .yaml files.
// question.yaml defines the debate topic and optional settings (rounds defaults to 5).
// All other .yaml files are parsed as agents (must have "name" and "model").
func (d *Demo) Load(dir string) error {
	questionData, err := os.ReadFile(filepath.Join(dir, "question.yaml"))
	if err != nil {
		return fmt.Errorf("reading question.yaml: %w", err)
	}

	var qf questionFile
	if err := yaml.Unmarshal(questionData, &qf); err != nil {
		return fmt.Errorf("parsing question.yaml: %w", err)
	}

	d.Question = strings.TrimSpace(qf.Question)
	if d.Question == "" {
		return fmt.Errorf("missing 'question' in question.yaml")
	}
	d.Rounds = qf.Rounds
	if d.Rounds == 0 {
		d.Rounds = 5
	} else if d.Rounds < 0 {
		return fmt.Errorf("invalid 'rounds': must be 0 or greater")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory: %w", err)
	}

	d.Agents = nil
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".yaml") || strings.EqualFold(name, "question.yaml") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("reading %s: %w", name, err)
		}

		var agent Agent
		if err := yaml.Unmarshal(data, &agent); err != nil {
			return fmt.Errorf("parsing %s: %w", name, err)
		}
		agent.Name = strings.TrimSpace(agent.Name)
		agent.Model = strings.TrimSpace(agent.Model)
		if agent.Name == "" {
			return fmt.Errorf("missing 'name' in %s", name)
		}
		if agent.Model == "" {
			return fmt.Errorf("missing 'model' in %s", name)
		}
		agent.Instructions = strings.TrimSpace(agent.Instructions)
		if err := validateAgentConfig(agent, name); err != nil {
			return err
		}
		d.Agents = append(d.Agents, agent)
	}

	sort.Slice(d.Agents, func(i, j int) bool {
		return d.Agents[i].Name < d.Agents[j].Name
	})

	return nil
}

func validateAgentConfig(agent Agent, filename string) error {
	if agent.MaxTokens < 0 {
		return fmt.Errorf("invalid 'max_tokens' in %s: must be 0 or greater", filename)
	}
	if err := validateOptionalRange("temperature", agent.Temperature, 0, 2, filename); err != nil {
		return err
	}
	if err := validateOptionalRange("top_p", agent.TopP, 0, 1, filename); err != nil {
		return err
	}
	return nil
}

func validateOptionalRange(field string, value *float64, min float64, max float64, filename string) error {
	if value == nil {
		return nil
	}
	if math.IsNaN(*value) || math.IsInf(*value, 0) || *value < min || *value > max {
		return fmt.Errorf("invalid %q in %s: must be between %g and %g", field, filename, min, max)
	}
	return nil
}
