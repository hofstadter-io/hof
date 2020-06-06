package cmdcontext

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/lib/config"
)

var getLong = `print a context or value(s) at path(s)`

func GetRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	if len(args) == 0 {
		val, err := config.GetRuntime().ContextGet("")
		if err != nil {
			return err
		}

		z := cue.Value{}
		if val == z {
			return fmt.Errorf("no context found, use 'hof context -h' to learn create and use contexts")
		}

		bytes, err := format.Node(val.Syntax())
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil
	}

	for _, a := range args {
		val, err := config.GetRuntime().ContextGet(a)
		if err != nil {
			return err
		}

		bytes, err := format.Node(val.Syntax())
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n\n", a, string(bytes))
	}

	return nil
}

var GetCmd = &cobra.Command{

	Use: "get <key.path>",

	Short: "print a context or value(s) at path(s)",

	Long: getLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

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

func init() {

	help := GetCmd.HelpFunc()
	usage := GetCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	GetCmd.SetHelpFunc(thelp)
	GetCmd.SetUsageFunc(tusage)

}