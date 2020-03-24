package studios

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/secret"
)

var secretLong = `Work with Hofstadter Studios secrets`

var SecretCmd = &cobra.Command{

	Use: "secret",

	Aliases: []string{
		"secrets",
		"shh",
	},

	Short: "Work with Hofstadter Studios secrets",

	Long: secretLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios secret")

	},
}

func init() {
	SecretCmd.AddCommand(secret.ListCmd)
	SecretCmd.AddCommand(secret.GetCmd)
	SecretCmd.AddCommand(secret.CreateCmd)
	SecretCmd.AddCommand(secret.UpdateCmd)
	SecretCmd.AddCommand(secret.DeleteCmd)
}
