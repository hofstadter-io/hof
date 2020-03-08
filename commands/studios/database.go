package studios

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
	DatabaseCmd.AddCommand(database.StatusCmd)
	DatabaseCmd.AddCommand(database.ResetCmd)
	DatabaseCmd.AddCommand(database.CheckpointCmd)
}
