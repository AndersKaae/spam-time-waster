package gemini

import (
	"context"
	"fmt"

	"github.com/anderskaae/spam-time-waster/internal/config"
	"google.golang.org/genai"
)

func Prompt(ctx context.Context, config *config.Config, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: config.GeminiAPIKey,
	})
	if err != nil {
		return "", err
	}

	result, err := client.Models.GenerateContent(ctx, config.GeminiModel, genai.Text(prompt), nil)
	if err != nil {
		return "", err
	}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%v", result.Candidates[0].Content.Parts[0]), nil
	}

	return "No response", nil
}
