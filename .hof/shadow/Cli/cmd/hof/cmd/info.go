package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

var infoLong = `print information about known resources`

func init() {

	InfoCmd.Flags().BoolVarP(&(flags.InfoFlags.Builtin), "builtin", "", false, "Only print builtin resources")
	InfoCmd.Flags().BoolVarP(&(flags.InfoFlags.Custom), "custom", "", false, "Only print custom resources")
	InfoCmd.Flags().BoolVarP(&(flags.InfoFlags.Local), "here", "", false, "Only print workspace resources")
}

func InfoRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info",

	Aliases: []string{
		"i",
	},

	Short: "print information about known resources",

	Long: infoLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = InfoRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := InfoCmd.HelpFunc()
	usage := InfoCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "", 0)
		return usage(cmd)
	}
	InfoCmd.SetHelpFunc(thelp)
	InfoCmd.SetUsageFunc(tusage)

}
