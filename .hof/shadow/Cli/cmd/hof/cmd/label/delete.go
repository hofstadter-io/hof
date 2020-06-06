package cmdlabel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var deleteLong = `delete labels from your workspace or system`

func DeleteRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var DeleteCmd = &cobra.Command{

	Use: "delete",

	Aliases: []string{
		"del",
		"remove",
	},

	Short: "delete labels from your workspace or system",

	Long: deleteLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = DeleteRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := DeleteCmd.HelpFunc()
	usage := DeleteCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	DeleteCmd.SetHelpFunc(thelp)
	DeleteCmd.SetUsageFunc(tusage)

}
