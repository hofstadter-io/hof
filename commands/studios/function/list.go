package function

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/function"
)

var ListLong = `List your functions`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your functions",

	Long: ListLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof function list:")

		err := function.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
