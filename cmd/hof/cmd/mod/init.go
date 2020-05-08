package cmdmod

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var initLong = `initialize a new module in the current directory`

func InitRun(lang string, module string) (err error) {

	err = lib.Init(lang, module)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var InitCmd = &cobra.Command{

	Use: "init <lang> <module>",

	Short: "initialize a new module in the current directory",

	Long: initLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

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
			fmt.Println("missing required argument: 'Module'")
			cmd.Usage()
			os.Exit(1)
		}

		var module string

		if 1 < len(args) {

			module = args[1]

		}

		err = InitRun(lang, module)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
