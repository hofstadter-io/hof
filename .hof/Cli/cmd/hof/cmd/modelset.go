package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/modelset"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var modelsetLong = `create, view, migrate, and understand your modelsets.`

var ModelsetCmd = &cobra.Command{

	Use: "modelset",

	Aliases: []string{
		"mset",
	},

	Short: "create, view, migrate, and understand your modelsets.",

	Long: modelsetLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {

	help := ModelsetCmd.HelpFunc()
	usage := ModelsetCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "<omit>", 0)
		return usage(cmd)
	}
	ModelsetCmd.SetHelpFunc(thelp)
	ModelsetCmd.SetUsageFunc(tusage)

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
