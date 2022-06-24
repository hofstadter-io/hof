package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

var renderLong = `hof render joins CUE with an extended Go base text/template system
  https://docs.hofstadter.io/code-generation/template-writing/

# Render a template
hof render data.cue -T template.txt
hof render data.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
hof render data.cue -T template.txt -P partial.txt

# The template flag
hof render data.cue ...

  # Multiple templates
  -T templateA.txt -T templateB.txt

  # Cuepath to select sub-input values
  -T 'templateA.txt;foo'
  -T 'templateB.txt;sub.val'

  # Writing to file
  -T 'templateA.txt;;a.txt'
  -T 'templateB.txt;sub.val;b.txt'

  # Templated output path 
  -T 'templateA.txt;;{{ .name | ToLower }}.txt'

  # Repeated templates when input is a list
  #   The template will be processed per item
  #   This also requires using a templated outpath
  -T 'template.txt;items;out/{{ .filepath }}.txt'
`

func init() {

	RenderCmd.Flags().StringSliceVarP(&(flags.RenderFlags.Template), "template", "T", nil, "Template mappings to render with data from entrypoint as: <filepath>;<?cuepath>;<?outpath>")
	RenderCmd.Flags().StringSliceVarP(&(flags.RenderFlags.Partial), "partial", "P", nil, "file globs to partial templates to register with the templates")
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
