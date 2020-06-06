package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/ops"
)

var jumpLong = `Jumps help you do things with fewer keystrokes.`

func JumpRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = ops.RunJumpFromArgs(args)

	return err
}

var JumpCmd = &cobra.Command{

	Use: "jump",

	Aliases: []string{
		"j",
		"leap",
	},

	Short: "Jumps help you do things with fewer keystrokes.",

	Long: jumpLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = JumpRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := JumpCmd.HelpFunc()
	usage := JumpCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	JumpCmd.SetHelpFunc(thelp)
	JumpCmd.SetUsageFunc(tusage)

}