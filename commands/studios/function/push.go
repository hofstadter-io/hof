package function

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/function"
)

var PushLong = `Uploads the local copy and makes it the latest copy in Studios`

var PushCmd = &cobra.Command{

	Use: "push <function path>",

	Short: "Send and make the latest version on Studios",

	Long: PushLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof function push: ")

		err := function.Push()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
