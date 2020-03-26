package database

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var restoreLong = `Restore a Studios database.`

var RestoreCmd = &cobra.Command{

	Use: "restore <name or id> <backup-name>",

	Short: "Restore a Studios database.",

	Long: restoreLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Ident'")
			cmd.Usage()
			os.Exit(1)
		}

		var ident string

		if 0 < len(args) {

			ident = args[0]

		}

		// Default body

		fmt.Println("hof studios database restore", ident)

	},
}
