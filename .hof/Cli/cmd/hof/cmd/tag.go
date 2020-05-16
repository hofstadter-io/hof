package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var tagLong = `Create, list, delete or verify a tag object signed with GPG`

func TagRun(args []string) (err error) {

	return err
}

var TagCmd = &cobra.Command{

	Use: "tag",

	Short: "Create, list, delete or verify a tag object signed with GPG",

	Long: tagLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = TagRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
