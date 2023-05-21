package cmd

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
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

	c.Value = c.Value.FillPath(cue.ParsePath("Args"), args)
	c.Value = c.Value.FillPath(cue.ParsePath("Files"), files)
	c.Value = c.Value.FillPath(cue.ParsePath("Question"), question)

	err = c.Value.Decode(c)
	if err != nil {
		err = cuetils.ExpandCueError(err)
		return err
	}

	fmt.Println("name:", c.Name)
	fmt.Println("model:", c.Model)
	fmt.Println("description:", c.Description)
	fmt.Println("messages:", c.Messages)



	
	// print gens
	return nil
}
