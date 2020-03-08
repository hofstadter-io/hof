package secret

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/secret"
)

var DeleteLong = `Delete a secret file`

var DeleteCmd = &cobra.Command{

	Use: "delete <name or id>",

	Short: "Delete a secret",

	Long: DeleteLong,

	Run: func(cmd *cobra.Command, args []string) {

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'name'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof secret delete:",
				name,
			)
		*/

		err := secret.Delete(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
