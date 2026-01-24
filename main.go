package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	copilot "github.com/github/copilot-sdk/go"
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

	// Create a session with system prompt and model configuration
	session, err := client.CreateSession(&copilot.SessionConfig{
		Model: "gpt-4.1",
		SystemMessage: &copilot.SystemMessageConfig{
			Mode:    "append",
			Content: agentDescription,
		},
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
		case "session.idle":
			// Signal completion when session is idle
			done <- true
		case "session.error":
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
