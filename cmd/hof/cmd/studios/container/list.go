package cmdcontainer

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listLong = `List your Studios containers`

func ListRun(args []string) (err error) {

	return err
}

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your containers",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ListRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
