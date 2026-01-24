package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/github/copilot-sdk/go"
)

//go:embed .github/agents/git-committer.agent.md
var agentDescription string

func main() {
	// Create a new Copilot client
	client := copilot.NewClient(nil)
	
	// Start the Copilot CLI server
	if err := client.Start(); err != nil {
		log.Fatalf("Failed to start Copilot client: %v", err)
	}
	defer client.Stop()

	// Parse the agent configuration from the embedded markdown
	agentConfig := parseAgentConfig(agentDescription)
	
	// Create a session with the custom agent
	session, err := client.CreateSession(&copilot.SessionConfig{
		CustomAgents: []copilot.CustomAgentConfig{agentConfig},
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Track completion
	done := make(chan bool)
	var output strings.Builder

	// Subscribe to session events
	session.On(func(event copilot.SessionEvent) {
		switch event.Type {
		case "assistant.message":
			// Print assistant messages
			if event.Data.Content != nil && *event.Data.Content != "" {
				output.WriteString(*event.Data.Content)
				output.WriteString("\n")
			}
		case "assistant.message_delta":
			// Print message deltas if streaming is enabled
			if event.Data.DeltaContent != nil && *event.Data.DeltaContent != "" {
				fmt.Print(*event.Data.DeltaContent)
			}
		case "done":
			// Signal completion
			done <- true
		case "error":
			// Print errors
			if event.Data.Message != nil && *event.Data.Message != "" {
				fmt.Fprintf(os.Stderr, "Error: %s\n", *event.Data.Message)
			}
			done <- true
		}
	})

	// Send the prompt to commit the currently staged files
	prompt := "commit the currently staged files"
	_, err = session.Send(copilot.MessageOptions{
		Prompt: prompt,
	})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	// Wait for completion
	<-done

	// Print final output if not streaming
	if output.Len() > 0 {
		fmt.Print(output.String())
	}
}

// parseAgentConfig extracts configuration from the agent markdown file
func parseAgentConfig(markdown string) copilot.CustomAgentConfig {
	config := copilot.CustomAgentConfig{
		Name:        "git-committer",
		DisplayName: "Git Committer",
		Description: "Expert at examining staged files and creating appropriate commit messages",
	}

	// Parse the markdown frontmatter and content
	lines := strings.Split(markdown, "\n")
	inFrontmatter := false
	frontmatterEnd := 0
	var tools []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Detect frontmatter
		if trimmed == "---" {
			if !inFrontmatter {
				inFrontmatter = true
			} else {
				frontmatterEnd = i
				break
			}
			continue
		}

		// Parse frontmatter fields
		if inFrontmatter {
			if strings.HasPrefix(trimmed, "name:") {
				config.Name = strings.TrimSpace(strings.TrimPrefix(trimmed, "name:"))
			} else if strings.HasPrefix(trimmed, "description:") {
				config.Description = strings.TrimSpace(strings.TrimPrefix(trimmed, "description:"))
			} else if strings.HasPrefix(trimmed, "- ") && len(tools) >= 0 {
				// Parse tools list
				tool := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
				tools = append(tools, tool)
			} else if strings.HasPrefix(trimmed, "infer:") {
				inferStr := strings.TrimSpace(strings.TrimPrefix(trimmed, "infer:"))
				infer := inferStr == "true"
				config.Infer = &infer
			}
		}
	}

	// Set tools if found
	if len(tools) > 0 {
		config.Tools = tools
	}

	// Extract the prompt from the markdown content (everything after frontmatter)
	if frontmatterEnd > 0 && frontmatterEnd < len(lines)-1 {
		config.Prompt = strings.TrimSpace(strings.Join(lines[frontmatterEnd+1:], "\n"))
	}

	return config
}
