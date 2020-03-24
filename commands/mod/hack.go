package mod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var hackLong = `dev command`

var HackCmd = &cobra.Command{

	Use: "hack",

	Hidden: true,

	Short: "dev command",

	Long: hackLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		err := lib.Hack("", args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
