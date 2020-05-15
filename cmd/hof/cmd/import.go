package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var importLong = `Import and create a data model from a multitude of sources such as
SQL, NoSQL, object storage, and buckets.`

func ImportRun(args []string) (err error) {

	return err
}

var ImportCmd = &cobra.Command{

	Use: "import",

	Short: "import and create a data model from a multitude of sources",

	Long: importLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ImportRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
