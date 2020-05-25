package cmdst

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/structural"
)

var mergeLong = `merge <new> onto <orig>, replacing values and adding new ones`

func MergeRun(orig string, update string, entrypoints []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = structural.RunMergeFromArgs(orig, update, entrypoints)

	return err
}

var MergeCmd = &cobra.Command{

	Use: "merge <orig> <update> [...entrypoints]",

	Short: "merge <new> onto <orig>, replacing values and adding new ones",

	Long: mergeLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'orig'")
			cmd.Usage()
			os.Exit(1)
		}

		var orig string

		if 0 < len(args) {

			orig = args[0]

		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'update'")
			cmd.Usage()
			os.Exit(1)
		}

		var update string

		if 1 < len(args) {

			update = args[1]

		}

		var entrypoints []string

		if 2 < len(args) {

			entrypoints = args[2:]

		}

		err = MergeRun(orig, update, entrypoints)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := MergeCmd.HelpFunc()
	usage := MergeCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "<omit>", 0)
		return usage(cmd)
	}
	MergeCmd.SetHelpFunc(thelp)
	MergeCmd.SetUsageFunc(tusage)

}
