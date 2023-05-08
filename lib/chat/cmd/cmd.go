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
	case "dm", "data", "datamodel":
		// fmt.Println("embed:", rest)	
		if len(args) == 0 {
			return fmt.Errorf("no input for datamodel provided")
		}
		return dmCall(rest, rflags, cflags)

	case "embed":
		// fmt.Println("embed:", rest)	
		if len(args) == 0 {
			return fmt.Errorf("no input for embedding provided")
		}
		return embedCall(rest, rflags, cflags)

	default:
		return chatCall(args, rflags, cflags)

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
}

func chatCall(args []string, rflags flags.RootPflagpole, cflags flags.ChatFlagpole) error {
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
	if cflags.Prompt != "" {
		bs, err := os.ReadFile(cflags.Prompt)	
		if err != nil {
			return err
		}
		s := string(bs)
		msg := openai.ChatCompletionMessage{
			Role: "system",
			Content: s,
		}
		msgs = append(msgs, msg)
	}
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

func dmCall(args []string, rflags flags.RootPflagpole, cflags flags.ChatFlagpole) error {
	lines := []string{}

	sysMsg := dmPretextString

	// construct inputs, we append the first file-like input
	// to the pretext as our model, and output the result to the same
	file := ""
	for _, arg := range args {
		if arg == "-" {
			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			s := string(b)
			if file == "" {
				file = "-"
				sysMsg += s
			} else {
				lines = append(lines, s)
			}
		} else if info, err := os.Lstat(arg); err == nil && !info.IsDir() {
			bs, err := os.ReadFile(info.Name())	
			if err != nil {
				return err
			}
			s := string(bs)
			if file == "" {
				file = info.Name()
				sysMsg += s
			} else {
				lines = append(lines, s)
			}
		} else {
			lines = append(lines, arg)
		}
	}

	// create the user message, by starting with the current or first datamodel
	usrMsg := ""
	if file == "" {
		// no file, probably the first iteration?
		file = "dm.cue"
		_, err := os.Stat(file)
		if err != nil {
			// fmt.Println(err)
			// not found, new dm most liklye
			usrMsg = dmStartingJSON
		} else {
			bs, err := os.ReadFile("dm.cue")	
			if err != nil {
				return err
			}
			usrMsg = string(bs)
		}
	}

	usrMsg += strings.Join(lines, "\n")

	if rflags.Verbosity > 0 {
		fmt.Println(sysMsg)
		fmt.Println(usrMsg)
		fmt.Printf("\nlength: %d\n\n", len(sysMsg) + len(usrMsg))
	}

	// make our chat messages
	msgs := []openai.ChatCompletionMessage {
		// system message
		{
			Role: "system",
			Content: sysMsg,
		},
		// user instructions
		{
			Role: "user",
			Content: usrMsg,
		},
	}

	// make the call
	resp, err := chat.OpenaiChat(msgs, "gpt-3.5-turbo")
	if err != nil {
		return err
	}

	//
	// fixes
	//
	// remove any triple ticks, they keep showing up despite the prompt...
	resp += "\n"
	resp = strings.Replace(resp, "```\n", "\n", -1)
	resp = strings.Replace(resp, "```json", "", -1)

	// add a new line for writing
	resp = strings.TrimSpace(resp) + "\n"

	// Print the final model
	fmt.Println(resp)
	// also write the file
	if file != "-" {
		if cflags.Outfile != "" {
			file = cflags.Outfile
		}
		err := os.WriteFile(file, []byte(resp), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func embedCall(args []string, rflags flags.RootPflagpole, cflags flags.ChatFlagpole) error {
		inputs := []string{}
		for _, R := range args {
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
			f := args[d.Index]
			D[f] = d.Embedding
		}

		bs, err := json.Marshal(D)
		if err != nil {
			return err
		}
		fmt.Println(string(bs))

		return nil
}
