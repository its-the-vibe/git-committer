package main

import (
	"testing"
)

func TestParseAgentConfig(t *testing.T) {
	markdown := `---
name: git-committer
description: Expert at examining staged files and creating appropriate commit messages
tools:
  - bash
  - view
  - grep
  - glob
infer: false
---

## Persona

You are an expert Git committer who specializes in examining staged changes and crafting clear, conventional commit messages.

## Your Task

When invoked, you should commit the staged files.
`

	config := parseAgentConfig(markdown)

	if config.Name != "git-committer" {
		t.Errorf("Expected name 'git-committer', got '%s'", config.Name)
	}

	if config.Description != "Expert at examining staged files and creating appropriate commit messages" {
		t.Errorf("Expected description to match, got '%s'", config.Description)
	}

	if len(config.Tools) != 4 {
		t.Errorf("Expected 4 tools, got %d", len(config.Tools))
	}

	expectedTools := []string{"bash", "view", "grep", "glob"}
	for i, tool := range expectedTools {
		if i >= len(config.Tools) || config.Tools[i] != tool {
			t.Errorf("Expected tool[%d] to be '%s', got '%s'", i, tool, config.Tools[i])
		}
	}

	if config.Infer == nil || *config.Infer != false {
		t.Errorf("Expected infer to be false")
	}

	if config.Prompt == "" {
		t.Error("Expected prompt to be non-empty")
	}

	if len(config.Prompt) < 50 {
		t.Errorf("Expected prompt to contain content from markdown, got length %d", len(config.Prompt))
	}
}
