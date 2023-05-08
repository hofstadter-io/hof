package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/chat"
	// "github.com/hofstadter-io/hof/lib/runtime"
)

func Run(args []string, rflags flags.RootPflagpole, cflags flags.ChatFlagpole) error {
	if len(args) == 0 {
		return fmt.Errorf("no input provided")
	}

	// fmt.Printf("lib/chat.Run: %v %v %v\n", args, rflags, cflags)

	cmd, rest := args[0], args[1:]
	switch cmd {
	case "embed":
		// fmt.Println("embed:", rest)	
		if len(rest) == 0 {
			return fmt.Errorf("no input for embedding provided")
		}

		inputs := []string{}
		for _, R := range rest {
			bs, err := os.ReadFile(R)	
			if err != nil {
				return err
			}
			s := string(bs)
			inputs = append(inputs, s)
		}

		var req openai.EmbeddingRequest
		apiKey := os.Getenv("OPENAI_API_KEY")
		client := openai.NewClient(apiKey)

		req.Model = openai.AdaEmbeddingV2
		req.Input = inputs

		ctx := context.Background()
		resp, err := client.CreateEmbeddings(ctx, req)
		if err != nil {
			return err
		}

		D := map[string]any{}
		for _, d := range resp.Data {
			f := rest[d.Index]
			D[f] = d.Embedding
		}

		bs, err := json.Marshal(D)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))

		return nil

	default:

		lines := []string{}

		for _, arg := range args {
			if arg == "-" {
				s, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				lines = append(lines, string(s))

			} else if info, err := os.Lstat(arg); err == nil && !info.IsDir() {
				bs, err := os.ReadFile(info.Name())	
				if err != nil {
					return err
				}
				s := string(bs)
				lines = append(lines, s)
			} else {
				lines = append(lines, arg)
			}

		}

		m := strings.Join(lines, "\n")
		msgs := make([]openai.ChatCompletionMessage,0)
		msg := openai.ChatCompletionMessage{
			Role: "user",
			Content: m,
		}
		msgs = append(msgs, msg)
		resp, err := chat.OpenaiChat(msgs, "gpt-3.5-turbo")
		fmt.Println(resp)
		if err != nil {
			return err
		}
		return nil
	}

	// load our cue, for future use
	/*
	R, err := runtime.New(extra, rflags)
	if err != nil {
		return err
	}
	err = R.Load()
	if err != nil {
		return err
	}
	*/

	// load code
	/*
	cbytes, err := os.ReadFile(jsonfile)
	if err != nil {
		return err
	}
	code := string(cbytes)

	// possibly load inst
	if strings.HasPrefix(inst, "./") {
		ibytes, err := os.ReadFile(inst)
		if err != nil {
			return err
		}
		inst = string(ibytes)
	}

	// make call
	resp, err := chat.OpenaiChat(code, inst, cflags.Model)
	if err != nil {
		return err
	}

	// write code
	fmt.Println(resp)
	err = os.WriteFile(jsonfile, []byte(resp), 0644)
	if err != nil {
		return err
	}
	*/

	return nil
}
