package functions

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listLong = `List your Studios functions`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your functions",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios functions list")

	},
}
