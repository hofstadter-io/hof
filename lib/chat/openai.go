package chat

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"

	// "github.com/hofstadter-io/hof/lib/templates"
)

var apiKey string

func init() {
	apiKey = os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// fmt.Println("OPENAI_API_KEY environment var is missing\nVisit https://platform.openai.com/account/api-keys to get one")
	}
}

func OpenaiChat(messages []openai.ChatCompletionMessage, model string) (string, error) {
	// override default model in interactive | chat mode
	if !(strings.HasPrefix(model, "gpt-3") || strings.HasPrefix(model, "gpt-4")) {
		return "", fmt.Errorf("using chat compatible model: %q", model, "\n")
	}

	client := openai.NewClient(apiKey)

	// initial req setup
	var req openai.ChatCompletionRequest

	req.Model = model
	req.N = 1
	req.MaxTokens = 2500
	req.Temperature = 0.042
	req.TopP = 0.69
	req.Messages = messages


	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	R := resp.Choices
	final := R[0].Message.Content

	return final, nil

	// add our message to the conversation
	/* TODO, once we are interactive
	msg = openai.ChatCompletionMessage{
		Role: "assistant",
		Content: final,
	}
	req.Messages = append(req.Messages, msg)
	*/
}

func OpenaiEmbedding(inputs []string) (string, error) {

	client := openai.NewClient(apiKey)

	// initial req setup
	var req openai.EmbeddingRequest

	req.Model = openai.AdaEmbeddingV2
	req.Input = inputs

	ctx := context.Background()
	resp, err := client.CreateEmbeddings(ctx, req)
	if err != nil {
		return "", err
	}
	D := resp.Data
	final := fmt.Sprint(D[0].Embedding)

	return final, nil

	// add our message to the conversation
	/* TODO, once we are interactive
	msg = openai.ChatCompletionMessage{
		Role: "assistant",
		Content: final,
	}
	req.Messages = append(req.Messages, msg)
	*/
}

