package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/model"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var modelLong = `create, view, migrate, and understand your data models.`

var ModelCmd = &cobra.Command{

	Use: "model",

	Aliases: []string{
		"models",
	},

	Short: "create, view, migrate, and understand your data models.",

	Long: modelLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {

	help := ModelCmd.HelpFunc()
	usage := ModelCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		return usage(cmd)
	}
	ModelCmd.SetHelpFunc(thelp)
	ModelCmd.SetUsageFunc(tusage)

	ModelCmd.AddCommand(cmdmodel.SetCmd)
	ModelCmd.AddCommand(cmdmodel.CreateCmd)
	ModelCmd.AddCommand(cmdmodel.ViewCmd)
	ModelCmd.AddCommand(cmdmodel.ListCmd)
	ModelCmd.AddCommand(cmdmodel.StatusCmd)
	ModelCmd.AddCommand(cmdmodel.GraphCmd)
	ModelCmd.AddCommand(cmdmodel.DiffCmd)
	ModelCmd.AddCommand(cmdmodel.MigrateCmd)
	ModelCmd.AddCommand(cmdmodel.TestCmd)
	ModelCmd.AddCommand(cmdmodel.DeleteCmd)

}
