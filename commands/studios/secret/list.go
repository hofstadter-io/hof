package secret

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/secret"
)

var ListLong = `List your Studios secrets`

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "List your secrets",

	Long: ListLong,

	Run: func(cmd *cobra.Command, args []string) {

		/*
			fmt.Println("hof secret list:")
		*/

		err := secret.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
