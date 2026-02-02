package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/genai"
)

func handleErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	location := os.Getenv("GOOGLE_CLOUD_LOCATION")

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  project,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	handleErr(err)

	for result, err := range client.Models.GenerateContentStream(
		ctx,
		"gemini-3-flash-preview",
		genai.Text("Why is the sky blue?"),
		&genai.GenerateContentConfig{
			MaxOutputTokens: 512,
		},
	) {
		if err != nil {
			handleErr(err)
		}
		fmt.Print(result.Candidates[0].Content.Parts[0].Text)
	}

	// // Call the GenerateContent method.
	// result, err := client.Models.GenerateContent(ctx,
	// 	"gemini-2.5-flash",
	// 	genai.Text("Tell me about New York?"),
	// 	&genai.GenerateContentConfig{
	// 		Temperature:      genai.Ptr[float32](0.5),
	// 		TopP:             genai.Ptr[float32](0.5),
	// 		TopK:             genai.Ptr[float32](2.0),
	// 		ResponseMIMEType: "application/json",
	// 		StopSequences:    []string{"\n"},
	// 		CandidateCount:   2,
	// 		Seed:             genai.Ptr[int32](42),
	// 		MaxOutputTokens:  128,
	// 		PresencePenalty:  genai.Ptr[float32](0.5),
	// 		FrequencyPenalty: genai.Ptr[float32](0.5),
	// 	},
	// )
	// handleErr(err)
	// debugPrint(result)

	// ms := client.Models.All(ctx)
	// for m := range ms {
	// 	fmt.Println(m.Name, m.Version)
	// }

}

func debugPrint[T any](r *T) {
	response, err := json.MarshalIndent(*r, "", "  ")
	handleErr(err)
	fmt.Println(string(response))
}
