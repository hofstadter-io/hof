package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	flowcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/middleware"
	"github.com/hofstadter-io/hof/flow/tasks"
	"github.com/hofstadter-io/hof/flow/flow"

	"github.com/hofstadter-io/hof/lib/chat"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

func With(name string, extra []string, rflags flags.RootPflagpole, cflags flags.ChatPflagpole) error {
	args := []string{}
	files := map[string]string{}
	entrypoints := []string{}
	var question string

	// parse incoming args
	for _, e := range extra {
		if strings.HasPrefix(e, ".") {
			entrypoints = append(entrypoints, e)

		} else if strings.HasPrefix(e, "@") {
			fn := e[1:]
			bs, err := os.ReadFile(fn)
			if err != nil {
				return err
			}
			files[fn] = string(bs)

		} else {
			parts := strings.Fields(e)
			if len(parts) == 1 {
				args = append(args, e)
			} else {
				question = e
			}
		}
	}

	// build our runtime up
	R, err := prepRuntime(entrypoints, rflags)
	if err != nil {
		return err
	}

	// find the chat entry the user wants
	var c *chat.Chat	
	for _, C := range R.Chats {
		if C.Hof.Metadata.Name == name {
			c = C
		}
	}
	if c == nil {
		return fmt.Errorf("no chat %q found", name)
	}

	// fill in values from user before decoding
	c.Value = c.Value.FillPath(cue.ParsePath("Args"), args)
	c.Value = c.Value.FillPath(cue.ParsePath("Files"), files)
	c.Value = c.Value.FillPath(cue.ParsePath("Question"), question)

	// get the latest decoded version of the value
	err = c.Value.Decode(c)
	if err != nil {
		err = cuetils.ExpandCueError(err)
		return err
	}

	// do this after filling incase there is a default or calculated value
	if len(c.Question) == 0 {
		return fmt.Errorf("no question provided")
	}

	// decide which model to call, expand short names
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

	// some debug printing
	if rflags.Verbosity > 0 {
		fmt.Println("name:      ", c.Name)
		fmt.Println("model:     ", model)
		fmt.Println("question:  ", c.Question)
		fmt.Println("messages:  ", c.Messages)
		fmt.Println("params:    ", c.Parameters)
	}

	// maybe run the pre flow
	preExec := c.Value.LookupPath(cue.ParsePath("PreExec"))
	if preExec.Exists() && preExec.IsConcrete() {
		if rflags.Verbosity > 0 {
			fmt.Println("running pre exec flow:", preExec)
		}

		ctx := flowcontext.New()
		ctx.RootValue = preExec
		ctx.Stdin = os.Stdin
		ctx.Stdout = os.Stdout
		ctx.Stderr = os.Stderr
		ctx.Verbosity = rflags.Verbosity

		// how to inject tags into original value
		// fill / return value
		middleware.UseDefaults(ctx, &rflags, &flags.FlowFlags)
		tasks.RegisterDefaults(ctx)

		p, err := flow.NewFlow(ctx, preExec)
		if err != nil {
			return err
		}

		err = p.Start()
		if err != nil {
			return err
		}

		c.Value = c.Value.FillPath(cue.ParsePath("PreExec"), preExec)
		if c.Value.Err() != nil {
			return err
		}

	} else if preExec.Err() != nil {
		return preExec.Err()
	}

	// call the model, get the body back
	var body string

	if isOpenai {
		body, err = chat.OpenaiChat(model, c.Messages, c.Parameters)
		if err != nil {
			return err
		}
	} else {
		body, err = chat.GoogleChat(model, c.Messages, c.Examples, c.Parameters)
		if err != nil {
			return err
		}
	}

	// unmarshal the JSON response and fill back in
	resp := map[string]any{}
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return err
	}
	c.Value = c.Value.FillPath(cue.ParsePath("Response"), resp)

	// get the latest decoded version of the value
	// getting a weird "value was rounded down error" from CUE
	//err = c.Value.Decode(c)
	//if err != nil {
	//  err = cuetils.ExpandCueError(err)
	//  return err
	//}

	postExec := c.Value.LookupPath(cue.ParsePath("PostExec"))
	if postExec.Exists() && postExec.IsConcrete() {
		if rflags.Verbosity > 0 {
			fmt.Println("running post exec flow:", postExec)
		}
		ctx := flowcontext.New()
		ctx.RootValue = postExec
		ctx.Stdin = os.Stdin
		ctx.Stdout = os.Stdout
		ctx.Stderr = os.Stderr
		ctx.Verbosity = rflags.Verbosity

		// how to inject tags into original value
		// fill / return value
		middleware.UseDefaults(ctx, &rflags, &flags.FlowFlags)
		tasks.RegisterDefaults(ctx)

		p, err := flow.NewFlow(ctx, postExec)
		if err != nil {
			return err
		}

		err = p.Start()
		if err != nil {
			return err
		}

		c.Value = c.Value.FillPath(cue.ParsePath("PostExec"), postExec)
		if c.Value.Err() != nil {
			return err
		}

	} else if postExec.Err() != nil {
		return postExec.Err()
	}


	outVal := c.Value.LookupPath(cue.ParsePath("Output"))
	out, err := outVal.String()
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(out)) > 0 {
		fmt.Println(out)
	} else {
		// we can skip reprocessing
		fmt.Println(body)
	}
	
	// print gens
	return nil
}
