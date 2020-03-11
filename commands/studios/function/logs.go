package function

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/function"
)

var LogsLong = `List the logs of your function`

var LogsCmd = &cobra.Command{

	Use: "logs",

	Short: "List the logs of your function",

	Long: LogsLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof function logs:")

		err := function.Logs()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
