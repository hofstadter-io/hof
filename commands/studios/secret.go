package studios

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/secret"
)

var SecretLong = `Work with your Studios Secrets`

var SecretCmd = &cobra.Command{

	Use: "secret",

	Aliases: []string{
		"secrets",
		"shh",
	},

	Short: "Work with your Studios Secrets",

	Long: SecretLong,
}

func init() {
	// add sub-commands
	SecretCmd.AddCommand(secret.ListCmd)
	SecretCmd.AddCommand(secret.CreateCmd)
	SecretCmd.AddCommand(secret.UpdateCmd)
	SecretCmd.AddCommand(secret.DeleteCmd)
}
