package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/mod"
)

var modLong = `Hof has mvs embedded, so you can do all the same things from this subcommand`

var ModCmd = &cobra.Command{

	Use: "mod",

	Aliases: []string{
		"m",
	},

	Short: "manage project modules",

	Long: modLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof mod")

	},
}

func init() {
	ModCmd.AddCommand(mod.InfoCmd)
	ModCmd.AddCommand(mod.ConvertCmd)
	ModCmd.AddCommand(mod.GraphCmd)
	ModCmd.AddCommand(mod.StatusCmd)
	ModCmd.AddCommand(mod.InitCmd)
	ModCmd.AddCommand(mod.TidyCmd)
	ModCmd.AddCommand(mod.VendorCmd)
	ModCmd.AddCommand(mod.VerifyCmd)
	ModCmd.AddCommand(mod.HackCmd)
}
