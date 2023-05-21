package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"

	// "github.com/hofstadter-io/hof/lib/templates"
)

var openaiApiKey string

func init() {
	openaiApiKey = os.Getenv("OPENAI_API_KEY")
	if openaiApiKey == "" {
		// fmt.Println("OPENAI_API_KEY environment var is missing\nVisit https://platform.openai.com/account/api-keys to get one")
	}
}

func OpenaiChat(model string, messages []Message, params map[string]any) (string, error) {
	// override default model in interactive | chat mode
	if !(strings.HasPrefix(model, "gpt-3") || strings.HasPrefix(model, "gpt-4")) {
		return "", fmt.Errorf("incompatible openai chat model: %q", model)
	}

	if openaiApiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment var is missing\nVisit https://platform.openai.com/account/api-keys to get one")
	}

	client := openai.NewClient(openaiApiKey)

	msgs := make([]openai.ChatCompletionMessage,0,len(messages))
	for _, M := range messages {
		msg := openai.ChatCompletionMessage{
			Role: M.Role,
			Content: M.Content,
		}
		msgs = append(msgs, msg)
	}
	// initial req setup
	var req openai.ChatCompletionRequest
	var err error

	req.Model = model
	req.Messages = msgs

	if I, ok := params["N"]; ok {
		i, ok := I.(int)
		if !ok {
			return "", fmt.Errorf("N should be an int for ChatGPT, got %q", I)
		}
		req.N = i
	}
	if I, ok := params["MaxTokens"]; ok {
		i, ok := I.(int)
		if !ok {
			return "", fmt.Errorf("MaxTokens should be an int for ChatGPT, got %q", I)
		}
		req.MaxTokens = i
	}
	if F, ok := params["Temperature"]; ok {
		f, ok := F.(float64)
		if !ok {
			return "", fmt.Errorf("Temperature should be an float for ChatGPT, got %q", F)
		}
		req.Temperature = float32(f)
	}
	if F, ok := params["TopP"]; ok {
		f, ok := F.(float64)
		if !ok {
			return "", fmt.Errorf("TopP should be an float for ChatGPT, got %q", F)
		}
		req.TopP = float32(f)
	}
	if S, ok := params["Stop"]; ok {
		s, ok := S.([]string)
		if !ok {
			return "", fmt.Errorf("Stop should be a string for ChatGPT, got %q", S)
		}
		if s != nil {
			req.Stop = s
		}
	}



	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	body, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func OpenaiExtractContent(body string) (string, error) {
	d := map[string]any{}
	err := json.Unmarshal([]byte(body), &d)
	if err != nil {
		return body, err
	}

	c, ok := d["choices"]
	if !ok {
		return body, fmt.Errorf("Error: no choices found")
	}
	C, ok := c.([]any)
	if !ok {
		return body, fmt.Errorf("Error: choices not a list")
	}
	C0, ok := C[0].(map[string]any)
	if !ok {
		return body, fmt.Errorf("Error: choice 0 not a map")
	}

	m, ok := C0["message"]
	if !ok {
		return body, fmt.Errorf("Error: choice 0 missing message")
	}
	M, ok := m.(map[string]any)
	if !ok {
		return body, fmt.Errorf("Error: message not a map")
	}

	content, ok := M["content"]
	if !ok {
		return body, fmt.Errorf("Error: no content found")
	}

	ret, ok := content.(string)
	if !ok {
		return body, fmt.Errorf("Error: content is not a string")
	}


	//R := resp.Choices
	//final := R[0].Message.Content

	return ret, nil
}

func OpenaiEmbedding(inputs []string) (string, error) {

	client := openai.NewClient(openaiApiKey)

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

