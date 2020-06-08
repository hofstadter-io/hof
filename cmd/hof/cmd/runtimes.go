package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/runtimes"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var runtimesLong = `work with runtimes (go, js, py, bash, docker, cloud-vms, k8s, custom)`

var RuntimesCmd = &cobra.Command{

	Use: "runtimes",

	Short: "work with runtimes (go, js, py, bash, docker, cloud-vms, k8s, custom)",

	Long: runtimesLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {

	help := RuntimesCmd.HelpFunc()
	usage := RuntimesCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	RuntimesCmd.SetHelpFunc(thelp)
	RuntimesCmd.SetUsageFunc(tusage)

	RuntimesCmd.AddCommand(cmdruntimes.InfoCmd)
	RuntimesCmd.AddCommand(cmdruntimes.CreateCmd)
	RuntimesCmd.AddCommand(cmdruntimes.GetCmd)
	RuntimesCmd.AddCommand(cmdruntimes.SetCmd)
	RuntimesCmd.AddCommand(cmdruntimes.EditCmd)
	RuntimesCmd.AddCommand(cmdruntimes.DeleteCmd)
	RuntimesCmd.AddCommand(cmdruntimes.InstallCmd)
	RuntimesCmd.AddCommand(cmdruntimes.UninstallCmd)

}
