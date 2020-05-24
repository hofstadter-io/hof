package cmdst

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var maskLong = `mask <what> Cue value(s) from <orig>, thereby 'filtering' the original`

func MaskRun(orig string, what string, entrypoints []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var MaskCmd = &cobra.Command{

	Use: "mask <orig> <what> [...entrypoints]",

	Short: "mask <what> Cue value(s) from <orig>, thereby 'filtering' the original",

	Long: maskLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Orig'")
			cmd.Usage()
			os.Exit(1)
		}

		var orig string

		if 0 < len(args) {

			orig = args[0]

		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'What'")
			cmd.Usage()
			os.Exit(1)
		}

		var what string

		if 1 < len(args) {

			what = args[1]

		}

		var entrypoints []string

		if 2 < len(args) {

			entrypoints = args[2:]

		}

		err = MaskRun(orig, what, entrypoints)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := MaskCmd.HelpFunc()
	usage := MaskCmd.UsageFunc()

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
	MaskCmd.SetHelpFunc(thelp)
	MaskCmd.SetUsageFunc(tusage)

}
