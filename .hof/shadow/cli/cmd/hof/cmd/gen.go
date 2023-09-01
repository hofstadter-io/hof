package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/gen"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var genLong = `hof unifies CUE with Go's text/template system and diff3
  create on-liners to generate any file from any 
  build reusable and modular generators
  edit and regenerate those files while keeping changes

# Render a template
  hof gen input.cue -T template.txt
  hof gen input.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
  hof gen input.cue -T template.txt -P partial.txt

# The template flag as code gen mappings
  hof gen input.cue -T ... -T ...

  # Generate multiple templates at once
  -T templateA.txt -T templateB.txt

  # Select a sub-input value by CUEpath
  -T templateA.txt:foo
  -T templateB.txt:sub.val

  # Choose a schema with @
  -T templateA.txt:foo@Foo
  -T templateB.txt:sub.val@schemas.val

  # Writing to file with = (semicolon)
  -T templateA.txt=a.txt
  -T templateB.txt:sub.val@schema=b.txt

  # Templated output path, braces need quotes
  -T templateA.txt:='{{ .name | lower }}.txt'

	# Data Files are created when no template
  -T :sub.val='{{ .name | lower }}.json'

  # Repeated templates are used when
  # 1. the output has a '[]' prefix
  # 2. the input is a list or array
  #   The template will be processed per entry
  #   This also requires using a templated outpath
  -T template.txt:items='[]out/{{ .filepath }}.txt'
  -T :items='[]out/{{ .filepath }}.yaml'

  # Output everything to a directory (out name is the same)
  -O out -T types.go -T handlers.go

  # Watch files and directories, doing full or Xcue-less reloads
  -W *.cue -X *.go -O out -T types.go -T handlers.go

# Turn any hof gen flags into a reusable generator module
  hof gen [entrypoints] flags... --as-module [name]
  hof gen [entrypoints] -G [name]

# Bootstrap a new generator module
  hof gen init github.com/hofstadter-io/demos

# List availabel generators
	hof gen list

# Learn about writing templates, with extra functions and helpers
  https://docs.hofstadter.io/code-generation/template-writing/

# Check the tests for complete examples
  https://github.com/hofstadter-io/hof/tree/_dev/test/render

# Compose code gen mappings into reusable modules with
  hof gen app.cue -G frontend -G backend -G migrations -T ...
  https://docs.hofstadter.io/first-example/`

func init() {

	flags.SetupGenFlags(GenCmd.Flags(), &(flags.GenFlags))

}

func GenRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var GenCmd = &cobra.Command{

	Use: "gen [files...]",

	Aliases: []string{
		"G",
	},

	Short: "CUE powered code generation",

	Long: genLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

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

		ga.SendCommandPath(cmd.CommandPath() + " help")

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

	thelp := func(cmd *cobra.Command, args []string) {
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		return usage(cmd)
	}
	GenCmd.SetHelpFunc(thelp)
	GenCmd.SetUsageFunc(tusage)

	GenCmd.AddCommand(cmdgen.InitCmd)
	GenCmd.AddCommand(cmdgen.InfoCmd)
	GenCmd.AddCommand(cmdgen.ListCmd)

}
