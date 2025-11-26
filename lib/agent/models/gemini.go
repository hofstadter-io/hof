package models

import (
	"context"

	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func Gemini(ctx context.Context, model string) (model.LLM, error) {
	return gemini.NewModel(ctx, model, &genai.ClientConfig{
		// Project:  "gen-lang-client-0911744172",
		// Location: "us-central1",
		// Backend:  genai.BackendVertexAI,
	})
}
