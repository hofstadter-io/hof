package commands

import (

	// custom imports

	// infered imports

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/db"
)

// Tool:   hof
// Name:   db
// Usage:  db
// Parent: hof

var DbLong = `Work with your Studios DB`

var DbCmd = &cobra.Command{

	Use: "db",

	Short: "Work with your Studios DB",

	Long: DbLong,
}

func init() {
	RootCmd.AddCommand(DbCmd)
}

func init() {
	// add sub-commands to this command when present

	DbCmd.AddCommand(db.StatusCmd)
	DbCmd.AddCommand(db.ResetCmd)
	DbCmd.AddCommand(db.CheckpointCmd)
}
