package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

var renderLong = `generate arbitrary files from data and CUE entrypoints

hof render -t template.go data.cue > file.go`

func init() {

	RenderCmd.Flags().StringSliceVarP(&(flags.RenderFlags.Template), "template", "t", nil, "Template mappings to render with data from entrypoint as: filepath|cuepath|outpath")
}

func RenderRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var RenderCmd = &cobra.Command{

	Use: "render [flags] [entrypoints...]",

	Aliases: []string{
		"R",
	},

	Short: "generate arbitrary files from data and CUE entrypoints",

	Long: renderLong,

	PreRun: func(cmd *cobra.Command, args []string) {

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RenderRun(args)
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

	ohelp := RenderCmd.HelpFunc()
	ousage := RenderCmd.UsageFunc()
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

	RenderCmd.SetHelpFunc(help)
	RenderCmd.SetUsageFunc(usage)

}
