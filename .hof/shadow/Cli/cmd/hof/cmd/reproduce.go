package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var reproduceLong = `Record, share, and replay reproducible environments and processes`

func ReproduceRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var ReproduceCmd = &cobra.Command{

	Use: "reproduce",

	Aliases: []string{
		"repro",
	},

	Short: "Record, share, and replay reproducible environments and processes",

	Long: reproduceLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ReproduceRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := ReproduceCmd.HelpFunc()
	usage := ReproduceCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	ReproduceCmd.SetHelpFunc(thelp)
	ReproduceCmd.SetUsageFunc(tusage)

}
