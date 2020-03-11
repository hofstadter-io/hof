package function

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/function"
)

var ShutdownLong = `Shutsdown the function <name>, while preserving code in Studios.`

var ShutdownCmd = &cobra.Command{

	Use: "shutdown",

	Short: "Shutsdown the function <name>",

	Long: ShutdownLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof function shutdown:",
				name,
			)
		*/

		err := function.Shutdown(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
