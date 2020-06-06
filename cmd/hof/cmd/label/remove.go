package cmdlabel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/labels"
)

var removeLong = `find and remove labels from resources`

func RemoveRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = labels.RunRemoveLabelFromArgs(args)

	return err
}

var RemoveCmd = &cobra.Command{

	Use: "remove",

	Aliases: []string{
		"r",
	},

	Short: "find and remove labels from resources",

	Long: removeLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RemoveRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := RemoveCmd.HelpFunc()
	usage := RemoveCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	RemoveCmd.SetHelpFunc(thelp)
	RemoveCmd.SetUsageFunc(tusage)

}