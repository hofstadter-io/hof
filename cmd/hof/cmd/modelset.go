package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/modelset"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var modelsetLong = `create, view, migrate, and understand your data models.`

var ModelsetCmd = &cobra.Command{

	Use: "modelset",

	Aliases: []string{
		"model",
		"m",
	},

	Short: "create, view, migrate, and understand your data models.",

	Long: modelsetLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {
	hf := ModelsetCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		hf(cmd, args)
	}
	ModelsetCmd.SetHelpFunc(f)
	ModelsetCmd.AddCommand(cmdmodelset.CreateCmd)
	ModelsetCmd.AddCommand(cmdmodelset.ViewCmd)
	ModelsetCmd.AddCommand(cmdmodelset.ListCmd)
	ModelsetCmd.AddCommand(cmdmodelset.StatusCmd)
	ModelsetCmd.AddCommand(cmdmodelset.GraphCmd)
	ModelsetCmd.AddCommand(cmdmodelset.DiffCmd)
	ModelsetCmd.AddCommand(cmdmodelset.MigrateCmd)
	ModelsetCmd.AddCommand(cmdmodelset.TestCmd)
	ModelsetCmd.AddCommand(cmdmodelset.DeleteCmd)
}
