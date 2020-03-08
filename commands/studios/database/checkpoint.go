package db

import (
	"fmt"
	"os"

	"github.com/hofstadter-io/hof/pkg/studios/db"
	"github.com/spf13/cobra"
)

var CheckpointLong = `Checkpoints the DB, making only the necessary changes to the schema.`

var CheckpointCmd = &cobra.Command{

	Use: "checkpoint",

	Aliases: []string{
		"migrate",
	},

	Short: "Checkpoint the DB schema",

	Long: CheckpointLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof db checkpoint:")

		err := db.Migrate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
