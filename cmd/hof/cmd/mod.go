package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var modLong = `Hof has mvs embedded, so you can do all the same things from this subcommand`

func ModRun(args []string) (err error) {

	return err
}

var ModCmd = &cobra.Command{

	Use: "mod",

	Aliases: []string{
		"m",
	},

	Short: "manage project modules",

	Long: modLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ModRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	hf := ModCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		hf(cmd, args)
	}
	ModCmd.SetHelpFunc(f)
	ModCmd.AddCommand(cmdmod.InfoCmd)
	ModCmd.AddCommand(cmdmod.ConvertCmd)
	ModCmd.AddCommand(cmdmod.GraphCmd)
	ModCmd.AddCommand(cmdmod.StatusCmd)
	ModCmd.AddCommand(cmdmod.InitCmd)
	ModCmd.AddCommand(cmdmod.TidyCmd)
	ModCmd.AddCommand(cmdmod.VendorCmd)
	ModCmd.AddCommand(cmdmod.VerifyCmd)
	ModCmd.AddCommand(cmdmod.HackCmd)
}
