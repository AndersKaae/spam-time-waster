package main

import (
	"context"
	"fmt"
	"github.com/anderskaae/spam-time-waster/config"
	"github.com/anderskaae/spam-time-waster/gemini"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	ctx := context.Background()
	resp, err := gemini.Prompt(ctx, cfg, "Write a clever one-liner about Gophers.")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Gemini:", resp)
}
