package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	var d Demo

	tests := []struct {
		name       string
		input      string
		wantFields map[string]string
		wantBody   string
		wantErr    bool
	}{
		{
			name:       "valid frontmatter with body",
			input:      "---\nname: Test\nmodel: gpt-4o\n---\nHello world",
			wantFields: map[string]string{"name": "Test", "model": "gpt-4o"},
			wantBody:   "Hello world",
		},
		{
			name:       "frontmatter only, no body",
			input:      "---\nkey: value\n---",
			wantFields: map[string]string{"key": "value"},
			wantBody:   "",
		},
		{
			name:       "whitespace around values",
			input:      "---\n  name :  Agent A  \n---\nBody text",
			wantFields: map[string]string{"name": "Agent A"},
			wantBody:   "Body text",
		},
		{
			name:    "missing opening delimiter",
			input:   "name: Test\n---\nBody",
			wantErr: true,
		},
		{
			name:    "missing closing delimiter",
			input:   "---\nname: Test\nBody",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, body, err := d.ParseFrontmatter(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
			for k, want := range tt.wantFields {
				if got := fields[k]; got != want {
					t.Errorf("fields[%q] = %q, want %q", k, got, want)
				}
			}
		})
	}
}

func TestParseAgentFile(t *testing.T) {
	var d Demo

	tests := []struct {
		name    string
		input   string
		want    Agent
		wantErr bool
	}{
		{
			name:  "valid agent",
			input: "---\nname: Agent A\nmodel: gpt-4o\n---\nYou are helpful.",
			want:  Agent{Name: "Agent A", Model: "gpt-4o", Instructions: "You are helpful."},
		},
		{
			name:    "missing name",
			input:   "---\nmodel: gpt-4o\n---\nInstructions",
			wantErr: true,
		},
		{
			name:    "missing model",
			input:   "---\nname: Agent A\n---\nInstructions",
			wantErr: true,
		},
		{
			name:    "invalid frontmatter",
			input:   "no frontmatter here",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := d.ParseAgentFile(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if agent != tt.want {
				t.Errorf("got %+v, want %+v", agent, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("Question.md", "What is Go?")
	writeFile("AgentA.md", "---\nname: Alpha\nmodel: gpt-4o\n---\nBe concise.")
	writeFile("AgentB.md", "---\nname: Beta\nmodel: gpt-5-mini\n---\nBe critical.")

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

	writeFile("Question.md", "---\nrounds: 3\n---\nTest topic")
	writeFile("AgentA.md", "---\nname: A\nmodel: gpt-4o\n---\nHi")
	writeFile("AgentB.md", "---\nname: B\nmodel: gpt-4o\n---\nHi")

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

	// No frontmatter — should default to 5 rounds
	writeFile("Question.md", "Plain question")
	writeFile("AgentA.md", "---\nname: A\nmodel: gpt-4o\n---\nHi")
	writeFile("AgentB.md", "---\nname: B\nmodel: gpt-4o\n---\nHi")

	var demo Demo
	if err := demo.Load(dir); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if demo.Rounds != 5 {
		t.Errorf("rounds = %d, want 5 (default)", demo.Rounds)
	}
}

func TestLoadInvalidRounds(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("Question.md", "---\nrounds: abc\n---\nTopic")
	writeFile("AgentA.md", "---\nname: A\nmodel: gpt-4o\n---\nHi")

	var demo Demo
	if err := demo.Load(dir); err == nil {
		t.Fatal("expected error for invalid rounds")
	}
}

func TestLoadMissingQuestion(t *testing.T) {
	dir := t.TempDir()

	var demo Demo
	if err := demo.Load(dir); err == nil {
		t.Fatal("expected error for missing Question.md")
	}
}

func TestLoadSkipsNonMdFiles(t *testing.T) {
	dir := t.TempDir()

	writeFile := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	writeFile("Question.md", "Topic")
	writeFile("AgentA.md", "---\nname: A\nmodel: gpt-4o\n---\nHi")
	writeFile("AgentB.md", "---\nname: B\nmodel: gpt-4o\n---\nHi")
	writeFile("notes.txt", "should be ignored")

	var demo Demo
	if err := demo.Load(dir); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(demo.Agents) != 2 {
		t.Errorf("got %d agents, want 2 (txt file should be ignored)", len(demo.Agents))
	}
}
