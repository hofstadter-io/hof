package container

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listLong = `List your Studios containers`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your containers",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios container list")

	},
}
