package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listLong = `List your Studios apps`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your apps",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios app list")

	},
}
