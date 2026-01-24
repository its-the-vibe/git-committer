package main

import (
	_ "embed"
	"fmt"
	"log"

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

	// Subscribe to session events to display streaming output
	session.On(func(event copilot.SessionEvent) {
		switch event.Type {
		case "assistant.message_delta":
			// Print message deltas if streaming is enabled
			if event.Data.DeltaContent != nil && *event.Data.DeltaContent != "" {
				fmt.Print(*event.Data.DeltaContent)
			}
		}
	})

	// Send the prompt to commit the currently staged files and wait for completion
	prompt := "commit the currently staged files"
	_, err = session.SendAndWait(copilot.MessageOptions{
		Prompt: prompt,
	}, 0) // Use default 60s timeout
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	// Output is already printed via the streaming event handler above
	fmt.Println() // Add a final newline for clean output
}
