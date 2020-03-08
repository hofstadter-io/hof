package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/util"
)

var TestLong = `Test the context for authenticated connectivity`

var TestCmd = &cobra.Command{

	Use: "test",

	Short: "Test a context configuration",

	Long: TestLong,

	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Println("hof config test:")

		err := util.TestConfigAuth()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
