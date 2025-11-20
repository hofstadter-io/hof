package models

import (
	"context"
	"os"

	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func Gemini(ctx context.Context, model string) (model.LLM, error) {
	return gemini.NewModel(ctx, model, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
}
