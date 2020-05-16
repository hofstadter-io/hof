package cmdsecret

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var getLong = `print a secret`

func GetRun(args []string) (err error) {

	return err
}

var GetCmd = &cobra.Command{

	Use: "get",

	Short: "print a secret",

	Long: getLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = GetRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
