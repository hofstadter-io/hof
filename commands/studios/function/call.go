package function

import (
	"fmt"

	"github.com/spf13/cobra"
)

var callLong = `Call your Studios function`

var CallCmd = &cobra.Command{

	Use: "call",

	Short: "Call a function",

	Long: callLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios function call")

	},
}
