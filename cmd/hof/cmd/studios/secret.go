package cmdstudios

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/secret"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
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

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

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
	hf := SecretCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		as := strings.Join(args, "/")
		ga.SendGaEvent(c+"/help", as, 0)
		hf(cmd, args)
	}
	SecretCmd.SetHelpFunc(f)
	SecretCmd.AddCommand(cmdsecret.ListCmd)
	SecretCmd.AddCommand(cmdsecret.GetCmd)
	SecretCmd.AddCommand(cmdsecret.CreateCmd)
	SecretCmd.AddCommand(cmdsecret.UpdateCmd)
	SecretCmd.AddCommand(cmdsecret.DeleteCmd)
}
