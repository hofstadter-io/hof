package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/chat"
	// "github.com/hofstadter-io/hof/lib/runtime"
)

func Run(args []string, rflags flags.RootPflagpole, cflags flags.ChatFlagpole) error {
	if len(args) == 0 {
		return fmt.Errorf("no input provided")
	}

	return chatCall(args, rflags, cflags)
}

func chatCall(args []string, rflags flags.RootPflagpole, cflags flags.ChatFlagpole) error {
	// figure out model details
	isOpenai := true
	model := cflags.Model
	if strings.HasPrefix(model, "chat-") || model == "bard" {
		isOpenai = false
	}
	switch model {
	case "bard":
		model = "chat-bison"
	case "gpt3", "gpt-3":
		model = "gpt-3.5-turbo"
	case "gpt4":
		model = "gpt-4"
	}

	// all messages
	msgs := make([]chat.Message,0)
	exas := make([]chat.Example,0)

	// build up the system message
	if cflags.System != nil {
		sys := ""
		for _, S := range cflags.System {
			_, err := os.Lstat(S)
			if err == nil {
				// need to stat here
				bs, err := os.ReadFile(cflags.System[0])	
				if err != nil {
					return err
				}
				s := string(bs)
				sys += s
			} else {
				sys += S + "\n"
			}
		}

		msg := chat.Message{
			Role: "system",
			Content: sys,
		}
		msgs = append(msgs, msg)
	}

	// build up examples
	if cflags.Examples != nil {
		for e, E := range cflags.Examples {
			_, err := os.Lstat(E)
			if err == nil {
				// need to stat here
				bs, err := os.ReadFile(cflags.System[0])	
				if err != nil {
					return err
				}
				// unmarshal json
				var exs []chat.Example
				err = json.Unmarshal(bs, &exs)
				if err != nil {
					return err
				}
				exas = append(exas, exs...)
			} else {
				var ex chat.Example
				ip := strings.Index(E, "<INPUT>:")
				op := strings.Index(E, "<OUTPUT>:")
				if ip == -1 || op == -1 {
					return fmt.Errorf("example %d did not have the correct format in: %q", e, E)
				}
				ip += len("<INPUT>:") +1
				ex.Input = E[ip:op]
				op += len("<OUTPUT>:") +1
				ex.Output = E[op:]

				exas = append(exas, ex)
			}
		}

		// openai hack for example
		if isOpenai {
			sys := ""
			if len(msgs) == 1 {
				msgs[0].Content = sys
			}
			sys += "\nExamples:\n\n"
			for _, E := range exas {
				sys += fmt.Sprintf("input: %s\noutput:%s\n\n", E.Input, E.Output)	
			}
			if len(msgs) == 1 {
				msgs[0].Content = sys
			} else {
				msgs = append(msgs, chat.Message{
					Role: "system",
					Content: sys,
				})
			}
		}
	}

	if cflags.Messages != nil {
		for m, M := range cflags.Messages {
			_, err := os.Lstat(M)
			if err == nil {
				// need to stat here
				bs, err := os.ReadFile(cflags.System[0])	
				if err != nil {
					return err
				}
				// unmarshal json
				var msg []chat.Message
				err = json.Unmarshal(bs, &msg)
				if err != nil {
					return err
				}
				msgs = append(msgs, msg...)
			} else {
				var msg chat.Message
				ip := strings.Index(M, ">")
				if ip == -1 {
					return fmt.Errorf("message %d did not have the correct format in: %q", m, M)
				}
				msg.Role = M[:ip]
				msg.Content = M[ip+1:]

				msgs = append(msgs, msg)
			}
		}
	}

	// build up the user question or input
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
	msg := chat.Message{
		Role: "user",
		Content: m,
	}
	msgs = append(msgs, msg)
	// end user message

	// parse params
	params := make(map[string]any)
	params["N"] = cflags.Choices
	params["MaxTokens"] = cflags.MaxTokens
	params["Temperature"] = cflags.Temperature
	params["TopP"] = cflags.TopP
	params["TopK"] = cflags.TopK
	params["Stop"] = cflags.Stop

	var resp string

	if isOpenai {
		body, err := chat.OpenaiChat(model, msgs, params)
		if err != nil {
			return err
		}
		resp, err = chat.OpenaiExtractContent(body)
		if err != nil {
			return err
		}
	} else {
		body, err := chat.GoogleChat(model, msgs, exas, params)
		if err != nil {
			return err
		}
		resp, err = chat.GoogleExtractContent(body)
		if err != nil {
			return err
		}
	}

	fmt.Println(resp)
	return nil
}

