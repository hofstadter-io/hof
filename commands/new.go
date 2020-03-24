package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newLong = `create a new project or subcomponent or files, depending on the context`

var NewCmd = &cobra.Command{

	Use: "new",

	Short: "create a new project or subcomponent or files",

	Long: newLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof new")

	},
}
