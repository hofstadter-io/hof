package commands

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/app"
)

var AppLong = `Work with your Studios App`

var AppCmd = &cobra.Command{

	Use: "app",

	Aliases: []string{
		"apps",
	},

	Short: "Work with your Studios App",

	Long: AppLong,
}

func init() {
	// add sub-commands
	AppCmd.AddCommand(app.StatusCmd)
	AppCmd.AddCommand(app.VersionCmd)
	AppCmd.AddCommand(app.AvailableVersionCmd)
	AppCmd.AddCommand(app.ListCmd)
	AppCmd.AddCommand(app.CreateCmd)
	AppCmd.AddCommand(app.DeleteCmd)
	AppCmd.AddCommand(app.UpdateCmd)
	AppCmd.AddCommand(app.ResetCmd)
	AppCmd.AddCommand(app.ValidateCmd)
	AppCmd.AddCommand(app.GenerateCmd)
	AppCmd.AddCommand(app.SecretsCmd)
	AppCmd.AddCommand(app.PullCmd)
	AppCmd.AddCommand(app.PushCmd)
	AppCmd.AddCommand(app.DeployCmd)
	AppCmd.AddCommand(app.ShutdownCmd)
}
