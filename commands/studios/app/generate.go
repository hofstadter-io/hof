package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/app"
)

var GenerateLong = `Validate your application or components of it`

var GenerateCmd = &cobra.Command{

	Use: "generate",

	Aliases: []string{
		"gen",
		"g",
	},

	Short: "Validate your application",

	Long: GenerateLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof app generate:")

		err := app.Generate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
