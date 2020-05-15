package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var initLong = `init the current directory for hof usage.`

func InitRun(args []string) (err error) {

	return err
}

var InitCmd = &cobra.Command{

	Use: "init",

	Short: "init the current directory for hof usage.",

	Long: initLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = InitRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
