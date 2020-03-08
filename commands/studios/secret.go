package commands

import (

	// custom imports

	// infered imports

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/secret"
)

// Tool:   hof
// Name:   secret
// Usage:  secret
// Parent: hof

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
	RootCmd.AddCommand(SecretCmd)
}

func init() {
	// add sub-commands to this command when present

	SecretCmd.AddCommand(secret.ListCmd)
	SecretCmd.AddCommand(secret.CreateCmd)
	SecretCmd.AddCommand(secret.UpdateCmd)
	SecretCmd.AddCommand(secret.DeleteCmd)
}
