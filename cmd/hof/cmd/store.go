package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/store"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var storeLong = `create, checkpoint, and migrate your storage engines`

var StoreCmd = &cobra.Command{

	Use: "store",

	Aliases: []string{
		"s",
	},

	Short: "create, checkpoint, and migrate your storage engines",

	Long: storeLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {
	hf := StoreCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		hf(cmd, args)
	}
	StoreCmd.SetHelpFunc(f)
	StoreCmd.AddCommand(cmdstore.RunCmd)
	StoreCmd.AddCommand(cmdstore.ConnCmd)
}
