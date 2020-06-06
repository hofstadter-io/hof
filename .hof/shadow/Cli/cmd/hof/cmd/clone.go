package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var cloneLong = `clone a workspace or repository into a new directory`

func CloneRun(module string, name string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var CloneCmd = &cobra.Command{

	Use: "clone",

	Short: "clone a workspace or repository into a new directory",

	Long: cloneLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'module'")
			cmd.Usage()
			os.Exit(1)
		}

		var module string

		if 0 < len(args) {

			module = args[0]

		}

		var name string

		if 1 < len(args) {

			name = args[1]

		}

		err = CloneRun(module, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := CloneCmd.HelpFunc()
	usage := CloneCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	CloneCmd.SetHelpFunc(thelp)
	CloneCmd.SetUsageFunc(tusage)

}
