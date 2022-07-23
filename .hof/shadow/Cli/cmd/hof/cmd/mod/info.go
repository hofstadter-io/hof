package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"
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
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := InfoCmd.HelpFunc()
	ousage := InfoCmd.UsageFunc()
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

	InfoCmd.SetHelpFunc(help)
	InfoCmd.SetUsageFunc(usage)

}
