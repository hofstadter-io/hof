package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var infoLong = `print info about languages and modders known to hof mod
	- no arg prints a list of known languages
	- an arg prints info about the language modder configuration that would be used`

func InfoRun(lang string) (err error) {

	msg, err := mod.LangInfo(lang)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(msg)

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info [language]",

	Short: "print info about languages and modders known to hof mod",

	Long: infoLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		var lang string

		if 0 < len(args) {

			lang = args[0]

		}

		err = InfoRun(lang)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := InfoCmd.HelpFunc()
	usage := InfoCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	InfoCmd.SetHelpFunc(thelp)
	InfoCmd.SetUsageFunc(tusage)

}
