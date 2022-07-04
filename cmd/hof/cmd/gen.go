package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/lib/gen"
)

var genLong = `hof gen joins CUE with Go's text/template system and diff3
  create on-liners to generate any file from any data
  build reusable and modular generators
  edit and regenerate those files while keeping changes

If no generator is specified, hof gen runs in adhoc mode.

# Render a template
hof gen data.cue -T template.txt
hof gen data.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
hof gen data.cue -T template.txt -P partial.txt

# The template flag as code gen mappings
hof gen data.cue ...

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
  hof gen app.cue -G frontend -G backend -G migrations
  https://docs.hofstadter.io/first-example/

# You can mix adhof with generators by using
# both the -G and -T/-P flags`

func init() {

	GenCmd.Flags().BoolVarP(&(flags.GenFlags.Stats), "stats", "s", false, "Print generator statistics")
	GenCmd.Flags().StringSliceVarP(&(flags.GenFlags.Generator), "generator", "G", nil, "Generators to run, default is all discovered")
	GenCmd.Flags().StringSliceVarP(&(flags.GenFlags.Template), "template", "T", nil, "Template mappings to render as '<filepath>;<?cuepath>;<?outpath>'")
	GenCmd.Flags().StringSliceVarP(&(flags.GenFlags.Partial), "partial", "P", nil, "file globs to partial templates to register with the templates")
	GenCmd.Flags().BoolVarP(&(flags.GenFlags.Diff3), "diff3", "D", false, "enable diff3 support for adhoc render, generators are configured in code")
}

func GenRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = gen.Gen(args, flags.RootPflags, flags.GenFlags)

	return err
}

var GenCmd = &cobra.Command{

	Use: "gen [files...]",

	Aliases: []string{
		"G",
	},

	Short: "create arbitrary files from data with templates and generators",

	Long: genLong,

	PreRun: func(cmd *cobra.Command, args []string) {

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = GenRun(args)
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

	ohelp := GenCmd.HelpFunc()
	ousage := GenCmd.UsageFunc()
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

	GenCmd.SetHelpFunc(help)
	GenCmd.SetUsageFunc(usage)

}