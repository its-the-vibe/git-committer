package main

import (
	"context"
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

	ctx := context.Background()

	// Start the Copilot CLI server
	if err := client.Start(ctx); err != nil {
		log.Fatalf("Failed to start Copilot client: %v", err)
	}
	defer client.Stop()

	// Create a session with system prompt and model configuration
	session, err := client.CreateSession(ctx, &copilot.SessionConfig{
		Model:     "gpt-4.1",
		Streaming: true,
		SystemMessage: &copilot.SystemMessageConfig{
			Mode:    "replace",
			Content: agentDescription,
		},
		OnPermissionRequest: copilot.PermissionHandler.ApproveAll,
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Disconnect()

	// Subscribe to session events to display streaming output
	session.On(func(event copilot.SessionEvent) {
		switch event.Type {
		case "assistant.message_delta":
			// Print message deltas if streaming is enabled
			if event.Data.DeltaContent != nil && *event.Data.DeltaContent != "" {
				fmt.Print(*event.Data.DeltaContent)
			}
		case "session.idle":
			fmt.Println()
		case "elicitation.requested":
			fmt.Println("elicitation.requested")
			if event.Data.Content != nil && *event.Data.Content != "" {
				fmt.Println(*event.Data.Content)
			}
		}
	})

	// Send the prompt to commit the currently staged files and wait for completion
	prompt := "commit the currently staged files"
	_, err = session.SendAndWait(ctx, copilot.MessageOptions{
		Prompt: prompt,
	})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}
