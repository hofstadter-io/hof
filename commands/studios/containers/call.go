package containers

import (
	"fmt"

	"github.com/spf13/cobra"
)

var callLong = `Call your Studios container`

var CallCmd = &cobra.Command{

	Use: "call",

	Short: "Call a container",

	Long: callLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios containers call")

	},
}
