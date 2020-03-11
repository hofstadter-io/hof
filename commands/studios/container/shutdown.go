package container

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/container"
)

var ShutdownLong = `Shutsdown a container by name or from the current directory`

var ShutdownCmd = &cobra.Command{

	Use: "shutdown [name]",

	Short: "Shutsdown a container",

	Long: ShutdownLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof containers shutdown:",
				name,
			)
		*/

		err := container.Shutdown(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
