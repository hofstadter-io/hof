package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listLong = `List your Studios configs`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your configs",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios config list")

	},
}
