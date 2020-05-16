package cmdsecret

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var createLong = `create secrets`

func CreateRun(args []string) (err error) {

	return err
}

var CreateCmd = &cobra.Command{

	Use: "create",

	Short: "create secrets",

	Long: createLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CreateRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
