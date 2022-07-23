package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func LogoRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var LogoCmd = &cobra.Command{

	Use: "_",

	Hidden: true,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = LogoRun(args)
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

	ohelp := LogoCmd.HelpFunc()
	ousage := LogoCmd.UsageFunc()
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

	LogoCmd.SetHelpFunc(help)
	LogoCmd.SetUsageFunc(usage)

}
