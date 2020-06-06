package cmdst

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var queryLong = `query for values matching an expr and/or attributes`

func QueryRun(orig string, expr string, entrypoints []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var QueryCmd = &cobra.Command{

	Use: "query <orig> <expr> [...entrypoints]",

	Short: "query for values matching an expr and/or attributes",

	Long: queryLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "", 0)

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
			fmt.Println("missing required argument: 'expr'")
			cmd.Usage()
			os.Exit(1)
		}

		var expr string

		if 1 < len(args) {

			expr = args[1]

		}

		var entrypoints []string

		if 2 < len(args) {

			entrypoints = args[2:]

		}

		err = QueryRun(orig, expr, entrypoints)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := QueryCmd.HelpFunc()
	usage := QueryCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "", 0)
		return usage(cmd)
	}
	QueryCmd.SetHelpFunc(thelp)
	QueryCmd.SetUsageFunc(tusage)

}
