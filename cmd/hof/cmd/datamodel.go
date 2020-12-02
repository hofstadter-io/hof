package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/datamodel"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var datamodelLong = `Data models are sets of models which are used in many hof processes and modules.

At their core, they represent the most abstract representation for objects and
their relations in your applications. The are extended and annotated to add
context fot their usage in different circumstances: (DB vs Server vs Client).

Beyond representing models in their current form, a history is maintained so that:
  - database migrations can be created and managed
  - servers can handle multiple model versions
  - clients can implement feature flags
Much of this is actually handled by code generators and must be implemented there.
Hof handles the core data model definitions, history, and snapshot creation.`

func init() {

	DatamodelCmd.PersistentFlags().StringSliceVarP(&(flags.DatamodelPflags.Datamodels), "datamodel", "D", nil, "Datamodels for the datamodel commands")
	DatamodelCmd.PersistentFlags().StringSliceVarP(&(flags.DatamodelPflags.Modelsets), "modelset", "M", nil, "Modelsets for the datamodel commands")
	DatamodelCmd.PersistentFlags().StringSliceVarP(&(flags.DatamodelPflags.Models), "model", "m", nil, "Models for the datamodel commands")
}

var DatamodelCmd = &cobra.Command{

	Use: "datamodel",

	Aliases: []string{
		"dmod",
		"dm",
	},

	Short: "create, view, diff, calculate / migrate, and manage your data models",

	Long: datamodelLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := DatamodelCmd.HelpFunc()
	ousage := DatamodelCmd.UsageFunc()
	help := func(cmd *cobra.Command, args []string) {
		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

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
	DatamodelCmd.AddCommand(cmddatamodel.EditCmd)
	DatamodelCmd.AddCommand(cmddatamodel.DeleteCmd)
	DatamodelCmd.AddCommand(cmddatamodel.StatusCmd)
	DatamodelCmd.AddCommand(cmddatamodel.VisualizeCmd)
	DatamodelCmd.AddCommand(cmddatamodel.DiffCmd)
	DatamodelCmd.AddCommand(cmddatamodel.HistoryCmd)
	DatamodelCmd.AddCommand(cmddatamodel.MigrateCmd)
	DatamodelCmd.AddCommand(cmddatamodel.ApplyCmd)

}
