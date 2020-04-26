package cmdstudios

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/secret"
)

var secretLong = `Work with Hofstadter Studios secrets`

func SecretRun(args []string) (err error) {

	return err
}

var SecretCmd = &cobra.Command{

	Use: "secret",

	Aliases: []string{
		"secrets",
		"shh",
	},

	Short: "Work with Hofstadter Studios secrets",

	Long: secretLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = SecretRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	SecretCmd.AddCommand(cmdsecret.ListCmd)
	SecretCmd.AddCommand(cmdsecret.GetCmd)
	SecretCmd.AddCommand(cmdsecret.CreateCmd)
	SecretCmd.AddCommand(cmdsecret.UpdateCmd)
	SecretCmd.AddCommand(cmdsecret.DeleteCmd)
}
