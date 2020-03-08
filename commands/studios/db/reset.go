package db

import (
	// "fmt"

	// custom imports

	// infered imports

	"os"

	"github.com/hofstadter-io/hof/lib/db"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   reset
// Usage:  reset
// Parent: db

var ResetLong = `Resets the DB to the latest checkpoint & schema, also adding seed data.`

var (
	ResetHardFlag bool
)

func init() {
	ResetCmd.Flags().BoolVarP(&ResetHardFlag, "hard", "H", false, "perform a hard database rest, squashing all migrations into one.")
	viper.BindPFlag("hard", ResetCmd.Flags().Lookup("hard"))

}

var ResetCmd = &cobra.Command{

	Use: "reset",

	Short: "Reset the DB to the latest schema",

	Long: ResetLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In resetCmd", "args", args)
		// Argument Parsing

		// fmt.Println("hof db reset:")

		err := db.Reset(ResetHardFlag)
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
