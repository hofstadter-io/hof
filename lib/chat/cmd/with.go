package cmd

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/chat"
)

func With(name string, entrypoints []string, rflags flags.RootPflagpole, cflags flags.ChatPflagpole) error {
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

	fmt.Println("name:", c.Name)
	fmt.Println("model:", c.Model)
	fmt.Println("description:", c.Description)
	
	// print gens
	return nil
}
