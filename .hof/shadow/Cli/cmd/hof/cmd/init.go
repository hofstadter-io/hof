package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initLong = `Create a new workspace with initial files and registers with the global context.

When the name matches the current directory, the workspace is created there,
otherwise a new directory with the name will be created.`

func InitRun(module string, name string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var InitCmd = &cobra.Command{

	Use: "init <module> [name]",

	Short: "create an empty workspace or initialize an existing directory to one",

	Long: initLong,

	PreRun: func(cmd *cobra.Command, args []string) {

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

		err = InitRun(module, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := InitCmd.HelpFunc()
	ousage := InitCmd.UsageFunc()
	help := func(cmd *cobra.Command, args []string) {
		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	InitCmd.SetHelpFunc(help)
	InitCmd.SetUsageFunc(usage)

}
