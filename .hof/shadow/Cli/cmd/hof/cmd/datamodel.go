package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/datamodel"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var datamodelLong = `create, view, diff, calculate / migrate, and manage your data models`

var DatamodelCmd = &cobra.Command{

	Use: "datamodel",

	Aliases: []string{
		"dmod",
	},

	Short: "create, view, diff, calculate / migrate, and manage your data models",

	Long: datamodelLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {

	help := DatamodelCmd.HelpFunc()
	usage := DatamodelCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	DatamodelCmd.SetHelpFunc(thelp)
	DatamodelCmd.SetUsageFunc(tusage)

	DatamodelCmd.AddCommand(cmddatamodel.CreateCmd)
	DatamodelCmd.AddCommand(cmddatamodel.GetCmd)
	DatamodelCmd.AddCommand(cmddatamodel.SetCmd)
	DatamodelCmd.AddCommand(cmddatamodel.EditCmd)
	DatamodelCmd.AddCommand(cmddatamodel.DeleteCmd)
	DatamodelCmd.AddCommand(cmddatamodel.StatusCmd)
	DatamodelCmd.AddCommand(cmddatamodel.VisualizeCmd)
	DatamodelCmd.AddCommand(cmddatamodel.DiffCmd)
	DatamodelCmd.AddCommand(cmddatamodel.HistoryCmd)
	DatamodelCmd.AddCommand(cmddatamodel.MigrateCmd)
	DatamodelCmd.AddCommand(cmddatamodel.ApplyCmd)

}
