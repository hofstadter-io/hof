package cmdmod

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var convertLong = `convert another package system to MVS.`

func ConvertRun(lang string, filename string) (err error) {

	err = mod.Convert(lang, filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var ConvertCmd = &cobra.Command{

	Use: "convert <lang> <file>",

	Short: "convert another package system to MVS.",

	Long: convertLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Lang'")
			cmd.Usage()
			os.Exit(1)
		}

		var lang string

		if 0 < len(args) {

			lang = args[0]

		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'Filename'")
			cmd.Usage()
			os.Exit(1)
		}

		var filename string

		if 1 < len(args) {

			filename = args[1]

		}

		err = ConvertRun(lang, filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := ConvertCmd.HelpFunc()
	usage := ConvertCmd.UsageFunc()

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
	ConvertCmd.SetHelpFunc(thelp)
	ConvertCmd.SetUsageFunc(tusage)

}
