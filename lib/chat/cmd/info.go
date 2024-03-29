package cmd

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/chat"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

func Info(name string, entrypoints []string, rflags flags.RootPflagpole) error {
	R, err := prepRuntime(entrypoints, rflags)
	if err != nil {
		return err
	}

	// TODO...
	// 1. use table printer
	// 2. move this command up, large blocks of this ought
	var c *chat.Chat	
	for _, C := range R.Chats {
		if C.Hof.Metadata.Name == name {
			c = C
		}
	}
	if c == nil {
		return fmt.Errorf("no chat %q found", name)
	}

	err = c.Value.Decode(c)
	if err != nil {
		err = cuetils.ExpandCueError(err)
		return err
	}

	fmt.Println("name:        ", c.Name)
	fmt.Println("model:       ", c.Model)
	fmt.Println("description: ", c.Description)
	
	// print gens
	return nil
}
