package mod

import (

	// hello... something might need to go here

	"github.com/spf13/cobra"

	"fmt"

	"os"

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
