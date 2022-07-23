package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func GebRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var GebCmd = &cobra.Command{

	Use: "_geb",

	Hidden: true,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = GebRun(args)
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

	ohelp := GebCmd.HelpFunc()
	ousage := GebCmd.UsageFunc()
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

	GebCmd.SetHelpFunc(help)
	GebCmd.SetUsageFunc(usage)

}
