package container

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/container"
)

var ListLong = `List your containers`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your containers",

	Long: ListLong,

	Run: func(cmd *cobra.Command, args []string) {

		/*
			fmt.Println("hof containers list:")
		*/

		err := container.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
