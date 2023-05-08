package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/chat"
	// "github.com/hofstadter-io/hof/lib/runtime"
)

func Run(jsonfile, inst string, extra []string, rflags flags.RootPflagpole, cflags flags.ChatFlagpole) error {
	fmt.Printf("lib/chat.Run.%s %v %v %v\n", jsonfile, extra, rflags, cflags)

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
	/*
	err = os.WriteFile(jsonfile, []byte(resp), 0644)
	if err != nil {
		return err
	}
	*/

	return nil
}
