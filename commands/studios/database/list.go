package database

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listLong = `List your Studios databases`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your databases",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios database list")

	},
}
