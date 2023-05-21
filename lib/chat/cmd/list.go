package cmd

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func List(args []string, rflags flags.RootPflagpole) error {
	R, err := prepRuntime(args, rflags)
	if err != nil {
		return err
	}

	// TODO...
	// 1. use table printer
	// 2. move this command up, large blocks of this ought
	chats := make([]string, 0, len(R.Chats))
	for _, C := range R.Chats {
		chats = append(chats, C.Hof.Metadata.Name)
	}
	if len(chats) == 0 {
		return fmt.Errorf("no chats found")
	}
	fmt.Printf("Available Chats\n  ")
	fmt.Println(strings.Join(chats, "\n  "))
	
	// print gens
	return nil
}
