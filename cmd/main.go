package main

import (
	"context"
	"fmt"
	"github.com/anderskaae/spam-time-waster/internal/config"
	"github.com/anderskaae/spam-time-waster/internal/gemini"
	"github.com/anderskaae/spam-time-waster/internal/gmail"
	"log"
)

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

	labelName := cfg.GmailLabel
	label, err := gmailSvc.GetLabelByName("me", labelName)
	if err != nil {
		log.Fatalf("Error searching for label: %v", err)
	}

	if label == nil {
		log.Fatalf("Label '%s' not found. Please create it in your Gmail account to continue.", labelName)
	}

	fmt.Printf("Found label '%s' (ID: %s)\n", labelName, label.Id)

	// List messages with the specific label
	r, err := gmailSvc.Users.Messages.List("me").LabelIds(label.Id).MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages for label: %v", err)
	}

	fmt.Printf("Found %d messages with label '%s'\n", len(r.Messages), labelName)

	resp, err := gemini.Prompt(ctx, cfg, "Write a clever one-liner about Gophers.")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Gemini:", resp)
}
