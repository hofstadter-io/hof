package cmdmod

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var infoLong = `  print info about languages and modders known to mvs
    - no arg prints a list of known languages
    - an arg prints info about the language modder configuration that would be used`

func InfoRun(lang string) (err error) {

	msg, err := lib.LangInfo(lang)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(msg)

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info [language]",

	Short: "print info about languages and modders known to mvs",

	Long: infoLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

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
