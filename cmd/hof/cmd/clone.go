package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var cloneLong = `Clone a Workspace into a new directory`

func CloneRun(args []string) (err error) {

	return err
}

var CloneCmd = &cobra.Command{

	Use: "clone",

	Short: "Clone a Workspace into a new directory",

	Long: cloneLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CloneRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
