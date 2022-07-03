package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/lib/gen"
)

var renderLong = `hof render joins CUE with Go's text/template system and diff3
  create on-liners to generate any file from any data
  edit and regenerate those files while keeping changes

# Render a template
hof render data.cue -T template.txt
hof render data.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
hof render data.cue -T template.txt -P partial.txt

# The template flag as code gen mappings
hof render data.cue ...

  # Generate multiple templates at once
  -T templateA.txt -T templateB.txt

  # Select a sub-input value by CUEpath
  -T 'templateA.txt:foo'
  -T 'templateB.txt:sub.val'

  # Choose a schema with @
  -T 'templateA.txt:foo@#foo'
  -T 'templateB.txt:sub.val@schemas.val'

  # Writing to file with ; (semicolon)
  -T 'templateA.txt;a.txt'
  -T 'templateB.txt:sub.val@schema;b.txt'

  # Templated output path 
  -T 'templateA.txt:;{{ .name | lower }}.txt'

  # Repeated templates are used when
  # 1. the output has a '[]' prefix
  # 2. the input is a list or array
  #   The template will be processed per entry
  #   This also requires using a templated outpath
  -T 'template.txt:items;[]out/{{ .filepath }}.txt'

# Learn about writing templates, with extra functions and helpers
  https://docs.hofstadter.io/code-generation/template-writing/

# Check the tests for complete examples
  https://github.com/hofstadter-io/hof/tree/_dev/test/render

# Compose code gen mappings into reusable modules with
  hof gen app.cue -g frontend -g backend -g migrations
  https://docs.hofstadter.io/first-example/`

func init() {

	RenderCmd.Flags().StringSliceVarP(&(flags.RenderFlags.Template), "template", "T", nil, "Template mappings to render with data from entrypoint as: <filepath>;<?cuepath>;<?outpath>")
	RenderCmd.Flags().StringSliceVarP(&(flags.RenderFlags.Partial), "partial", "P", nil, "file globs to partial templates to register with the templates")
	RenderCmd.Flags().BoolVarP(&(flags.RenderFlags.Diff3), "diff3", "D", false, "enable diff3 support, requires the .hof shadow directory")
}

func RenderRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = gen.Render(args, flags.RootPflags, flags.RenderFlags)

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