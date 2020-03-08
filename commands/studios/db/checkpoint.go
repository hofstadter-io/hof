package db

import (
	// "fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/db"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   checkpoint
// Usage:  checkpoint
// Parent: db

var CheckpointLong = `Checkpoints the DB, making only the necessary changes to the schema.`

var CheckpointCmd = &cobra.Command{

	Use: "checkpoint",

	Aliases: []string{
		"migrate",
	},

	Short: "Checkpoint the DB schema",

	Long: CheckpointLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In checkpointCmd", "args", args)
		// Argument Parsing

		// fmt.Println("hof db checkpoint:")

		err := db.Migrate()
		if err != nil {
		 os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
