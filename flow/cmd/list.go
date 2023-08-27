package cmd

import (
	"fmt"
	"sort"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func List(args []string, rflags flags.RootPflagpole, cflags flags.FlowPflagpole) error {
	R, err := prepRuntime(args, rflags, cflags)
	if err != nil {
		return err
	}

	fmt.Println("Available Generators")
	flows := make([]string, 0, len(R.Workflows))
	for _, G := range R.Workflows {
		fmt.Println(" ", G.Hof.Flow.Name)
		flows = append(flows, G.Hof.Flow.Name)
	}
	sort.Strings(flows)

	// TODO...
	// 1. use table printer
	// 2. move this command up, large blocks of this ought
	//flows := make([]string, 0, len(R.Workflows))
	//for _, G := range R.Workflows {
	//  flows = append(flows, G.Hof.Flow.Name)
	//}
	//sort.Strings(flows)
	//fmt.Printf("Available Generators\n  ")
	//fmt.Println(strings.Join(flows, "\n  "))
	
	// print gens
	return nil
}
