package main

import (
	"context"
	"fmt"
	"github.com/anderskaae/spam-time-waster/internal/config"
	"github.com/anderskaae/spam-time-waster/internal/gemini"
	"github.com/anderskaae/spam-time-waster/internal/gmail"
	gmail_v1 "google.golang.org/api/gmail/v1"
	"log"
)

func lookForLabel(cfg *config.Config, gmailSvc *gmail.Service) *gmail_v1.Label {
	labelName := cfg.GmailLabel

	label, err := gmailSvc.GetLabelByName("me", labelName)
	if err != nil {
		log.Fatalf("Error searching for label: %v", err)
	}

	if label == nil {
		log.Fatalf("Label '%s' not found. Please create it in your Gmail account to continue.", labelName)
	}

	fmt.Printf("Found label '%s' (ID: %s)\n", labelName, label.Id)
	return label
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	ctx := context.Background()

	// Initialize Gmail Service
	gmailSvc, err := gmail.NewService(ctx, cfg.GmailClientID, cfg.GmailClientSecret, cfg.GmailTokenFile)
	if err != nil {
		log.Fatalf("Unable to initialize Gmail service: %v", err)
	}

	label := lookForLabel(cfg, gmailSvc)

	// List messages with the specific label and map them to our Message objects
	messages, err := gmailSvc.GetMessagesByLabel("me", label.Id)
	if err != nil {
		log.Fatalf("Unable to retrieve messages for label: %v", err)
	}

	fmt.Printf("Retrieved %d messages with label '%s':\n", len(messages), cfg.GmailLabel)
	for _, msg := range messages {
		fmt.Printf("--- Message ID: %s ---\n", msg.ID)
		fmt.Printf("--- Thread ID: %s ---\n", msg.ThreadID)

		fmt.Printf("From:    %s\n", msg.From)
		fmt.Printf("Subject: %s\n", msg.Subject)
		// We could print the body but it might be long
		fmt.Printf("Body:    %s...\n", msg.Body[:100])

		resp, err := gemini.Prompt(ctx, cfg,
			"Act as a interested but pedantic and slightly confused potential client. "+
				"Respond (only email body) to the unsolicited email below. "+
				"Constraint: Do not commit to calls, demos, or payments. "+
				"Goal: Ask for deep clarification on a trivial part of their pitch. "+
				"Style: Polite but circular. Keep the length proportional to their message. "+
				"Example: If they mention 'AI-driven growth,' ask for a definition of 'growth' in their specific accounting framework. "+
				"Sign off as 'Aiden'. No placeholders. "+
				"Email body: "+msg.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	}
}
