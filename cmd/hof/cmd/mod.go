package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	cuecmd "cuelang.org/go/cmd/cue/cmd"
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var modLong = `CUE module dependency management, imported from upstream CUE project`

func ModRun(args []string) (err error) {

	runCueCmd(append([]string{"mod"}, args...))

	return err
}

var ModCmd = &cobra.Command{

	Use: "mod",

	Short: "CUE module dependency management",

	Long: modLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = ModRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var modsubs = []string{
	"edit",
	"fix",
	"get",
	"init",
	"publish",
	"registry",
	"resolve",
	"tidy",
}

func init() {
	ModCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		// fmt.Println("HELP!:", cmd.CommandPath())
		ga.SendCommandPath(cmd.CommandPath() + " help")
		runCueCmd([]string{"mod", "--help"})
	})

	for _, sub := range modsubs {
		cmd := &cobra.Command{
			Use: sub,
			Run: func(cmd *cobra.Command, args []string) {
				ga.SendCommandPath(cmd.CommandPath())
				runCueCmd(os.Args[1:])
			},
			FParseErrWhitelist: cobra.FParseErrWhitelist{
				UnknownFlags: true,
			},
		}
		cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			ga.SendCommandPath(cmd.CommandPath() + " help")
			runCueCmd([]string{"mod", sub, "--help"})
		})

		ModCmd.AddCommand(cmd)
	}

}

func runCueCmd(args []string) {
	c, _ := cuecmd.New(args)

	var buf bytes.Buffer
	c.SetOutput(&buf)

	err := c.Run(context.Background())

	// todo, use copy / pipe / replace so we can stream output
	s := buf.String()
	s = strings.Replace(s, "cue ", "hof ", -1)
	fmt.Println(s)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
