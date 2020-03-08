package commands

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/database"
)

var DatabaseLong = `Work with your Studios DB`

var DatabaseCmd = &cobra.Command{

	Use: "db",

	Short: "Work with your Studios DB",

	Long: DatabaseLong,
}

func init() {
	// add sub-commands
	DatabaseCmd.AddCommand(db.StatusCmd)
	DatabaseCmd.AddCommand(db.ResetCmd)
	DatabaseCmd.AddCommand(db.CheckpointCmd)
}
