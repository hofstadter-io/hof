package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var newLong = `create a new project or subcomponent or files, depending on the context`

func NewRun(args []string) (err error) {

	return err
}

var NewCmd = &cobra.Command{

	Use: "new",

	Short: "create a new project or subcomponent or files",

	Long: newLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = NewRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
