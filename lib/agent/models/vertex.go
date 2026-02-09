package models

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func Vertex(ctx context.Context, model string) (model.LLM, error) {

	use := os.Getenv("GOOGLE_GENAI_USE_VERTEXAI")
	if use != "" { // todo, be more truthy
		// creds file
		// creds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
		loc := os.Getenv("GOOGLE_CLOUD_LOCATION")
		if loc == "" {
			loc = "global"
		}

		fmt.Println("vertex:", proj, loc)

		// default inference (same as Go SDK) (typically a service account)
		return gemini.NewModel(ctx, model, &genai.ClientConfig{
			Project:  proj,
			Location: loc,
			Backend:  genai.BackendVertexAI,
		})
	}

	// default inference (same as Go SDK)
	return gemini.NewModel(ctx, model, &genai.ClientConfig{})
}
