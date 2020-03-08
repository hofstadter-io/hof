package container

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/crun"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   push
// Usage:  push [name]
// Parent: container

var PushLong = `Uploads the local copy and makes it the latest copy in Studios`

var PushCmd = &cobra.Command{

	Use: "push [name]",

	Short: "Send the latest version on Studios",

	Long: PushLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In pushCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		/*
			fmt.Println("hof containers push:",
				name,
			)
		*/

		err := crun.Push(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
