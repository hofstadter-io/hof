package chat

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/hofstadter-io/hof/lib/templates"
)

func OpenaiChat(code, inst, model string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment var is missing\nVisit https://platform.openai.com/account/api-keys to get one")
	}

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
	req.Temperature = 0.0
	req.TopP = 1.0
	req.Messages = make([]openai.ChatCompletionMessage,0)

	// The prompt engineering message for the System
	// TBD if separating this out is beneficial, but docs seem to imply this
	sys := openai.ChatCompletionMessage{
		Role: "system",
		Content: pretextString,
	}
	req.Messages = append(req.Messages, sys)

	// The user prompt string
	t, err := templates.CreateFromString("prompt", promptTemplate, nil)
	if err != nil {
		return "", nil
	}
	data := map[string]string{
		"code": strings.TrimSpace(code),
		"inst": strings.TrimSpace(inst),
	}
	bs, err := t.Render(data)
	if err != nil {
		return "", nil
	}
	prompt := string(bs)

	msg := openai.ChatCompletionMessage{
		Role: "user",
		Content: prompt,
	}
	req.Messages = append(req.Messages, msg)

	// debug
	fmt.Printf("%s\n\n%s\n\n", pretextString, prompt)

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

