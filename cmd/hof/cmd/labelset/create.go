package cmdlabelset

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/labels"
)

var createLong = `add labelsets to your workspace or system`

func CreateRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = labels.RunCreateLabelsetFromArgs(args)

	return err
}

var CreateCmd = &cobra.Command{

	Use: "create",

	Aliases: []string{
		"c",
	},

	Short: "add labelsets to your workspace or system",

	Long: createLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CreateRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := CreateCmd.HelpFunc()
	usage := CreateCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	CreateCmd.SetHelpFunc(thelp)
	CreateCmd.SetUsageFunc(tusage)

}