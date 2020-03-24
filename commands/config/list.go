package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listLong = `list configurations`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "list configurations",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof config list")

	},
}
