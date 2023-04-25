package cmd

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func List(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole) error {
	R, err := prepRuntime(args, rflags, gflags)
	if err != nil {
		return err
	}

	// TODO...
	// 1. use table printer
	// 2. move this command up, large blocks of this ought
	gens := make([]string, 0, len(R.Generators))
	for _, G := range R.Generators {
		gens = append(gens, G.Hof.Metadata.Name)
	}
	if len(gens) == 0 {
		return fmt.Errorf("no generators found")
	}
	fmt.Printf("Available Generators\n  ")
	fmt.Println(strings.Join(gens, "\n  "))
	
	// print gens
	return nil
}
