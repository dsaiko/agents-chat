package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("question.yaml", "question: What is Go?\n")
	writeFile("AgentA.yaml", "name: Alpha\nmodel: gpt-4o\ninstructions: Be concise.\n")
	writeFile("AgentB.yaml", "name: Beta\nmodel: gpt-5-mini\ninstructions: Be critical.\n")

	var demo Demo
	if err := demo.Load(dir); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if demo.Question != "What is Go?" {
		t.Errorf("question = %q, want %q", demo.Question, "What is Go?")
	}

	if len(demo.Agents) != 2 {
		t.Fatalf("got %d agents, want 2", len(demo.Agents))
	}

	// Agents should be sorted by name
	if demo.Agents[0].Name != "Alpha" {
		t.Errorf("first agent = %q, want %q", demo.Agents[0].Name, "Alpha")
	}
	if demo.Agents[0].Model != "gpt-4o" {
		t.Errorf("first agent model = %q, want %q", demo.Agents[0].Model, "gpt-4o")
	}
	if demo.Agents[1].Name != "Beta" {
		t.Errorf("second agent = %q, want %q", demo.Agents[1].Name, "Beta")
	}
	if demo.Agents[1].Model != "gpt-5-mini" {
		t.Errorf("second agent model = %q, want %q", demo.Agents[1].Model, "gpt-5-mini")
	}
}

func TestLoadWithRounds(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("question.yaml", "rounds: 3\nquestion: Test topic\n")
	writeFile("AgentA.yaml", "name: A\nmodel: gpt-4o\ninstructions: Hi\n")
	writeFile("AgentB.yaml", "name: B\nmodel: gpt-4o\ninstructions: Hi\n")

	var demo Demo
	if err := demo.Load(dir); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if demo.Rounds != 3 {
		t.Errorf("rounds = %d, want 3", demo.Rounds)
	}
	if demo.Question != "Test topic" {
		t.Errorf("question = %q, want %q", demo.Question, "Test topic")
	}
}

func TestLoadDefaultRounds(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("question.yaml", "question: Plain question\n")
	writeFile("AgentA.yaml", "name: A\nmodel: gpt-4o\ninstructions: Hi\n")
	writeFile("AgentB.yaml", "name: B\nmodel: gpt-4o\ninstructions: Hi\n")

	var demo Demo
	if err := demo.Load(dir); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if demo.Rounds != 5 {
		t.Errorf("rounds = %d, want 5 (default)", demo.Rounds)
	}
}

func TestLoadWithAllAgentParams(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("question.yaml", "question: Test\n")
	writeFile("AgentA.yaml", "name: A\nmodel: gpt-4o\nmax_tokens: 2048\ntemperature: 0.9\ntop_p: 0.95\ninstructions: Be helpful.\n")

	var demo Demo
	if err := demo.Load(dir); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	agent := demo.Agents[0]
	if agent.MaxTokens != 2048 {
		t.Errorf("max_tokens = %d, want 2048", agent.MaxTokens)
	}
	if agent.Temperature != 0.9 {
		t.Errorf("temperature = %f, want 0.9", agent.Temperature)
	}
	if agent.TopP != 0.95 {
		t.Errorf("top_p = %f, want 0.95", agent.TopP)
	}
}

func TestLoadMissingName(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("question.yaml", "question: Topic\n")
	writeFile("AgentA.yaml", "model: gpt-4o\ninstructions: Hi\n")

	var demo Demo
	if err := demo.Load(dir); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestLoadMissingModel(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("question.yaml", "question: Topic\n")
	writeFile("AgentA.yaml", "name: A\ninstructions: Hi\n")

	var demo Demo
	if err := demo.Load(dir); err == nil {
		t.Fatal("expected error for missing model")
	}
}

func TestLoadMissingQuestion(t *testing.T) {
	dir := t.TempDir()

	var demo Demo
	if err := demo.Load(dir); err == nil {
		t.Fatal("expected error for missing question.yaml")
	}
}

func TestLoadSkipsNonYamlFiles(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("question.yaml", "question: Topic\n")
	writeFile("AgentA.yaml", "name: A\nmodel: gpt-4o\ninstructions: Hi\n")
	writeFile("AgentB.yaml", "name: B\nmodel: gpt-4o\ninstructions: Hi\n")
	writeFile("notes.txt", "should be ignored")
	writeFile("readme.md", "should also be ignored")

	var demo Demo
	if err := demo.Load(dir); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(demo.Agents) != 2 {
		t.Errorf("got %d agents, want 2 (non-yaml files should be ignored)", len(demo.Agents))
	}
}
