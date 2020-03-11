package function

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/function"
)

var VersionsLong = `Get the supported runtime versions for Hofstadter Studios`

var VersionsCmd = &cobra.Command{

	Use: "versions",

	Short: "Get the runtime versions",

	Long: VersionsLong,

	Run: func(cmd *cobra.Command, args []string) {

		/*
			fmt.Println("hof function versions:")
		*/

		err := function.Versions()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
