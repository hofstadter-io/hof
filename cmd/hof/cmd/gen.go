package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/gen"
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
  hof gen --init github.com/hofstadter-io/demos

# Learn about writing templates, with extra functions and helpers
  https://docs.hofstadter.io/code-generation/template-writing/

# Check the tests for complete examples
  https://github.com/hofstadter-io/hof/tree/_dev/test/render

# Compose code gen mappings into reusable modules with
  hof gen app.cue -G frontend -G backend -G migrations -T ...
  https://docs.hofstadter.io/first-example/`

func init() {

	GenCmd.Flags().BoolVarP(&(flags.GenFlags.List), "list", "l", false, "list available generators")
	GenCmd.Flags().BoolVarP(&(flags.GenFlags.Stats), "stats", "s", false, "print generator statistics")
	GenCmd.Flags().StringSliceVarP(&(flags.GenFlags.Generator), "generator", "G", nil, "generator tags to run, default is all, or none if -T is used")
	GenCmd.Flags().StringSliceVarP(&(flags.GenFlags.Template), "template", "T", nil, "template mapping to render, see help for format")
	GenCmd.Flags().StringSliceVarP(&(flags.GenFlags.Partial), "partial", "P", nil, "file globs to partial templates to register with the templates")
	GenCmd.Flags().BoolVarP(&(flags.GenFlags.Diff3), "diff3", "D", false, "enable diff3 support for custom code")
	GenCmd.Flags().BoolVarP(&(flags.GenFlags.NoFormat), "no-format", "", false, "disable formatting during code gen (ad-hoc only)")
	GenCmd.Flags().BoolVarP(&(flags.GenFlags.Watch), "watch", "w", false, "run in watch mode, regenerating when files change, implied by -W/X")
	GenCmd.Flags().StringSliceVarP(&(flags.GenFlags.WatchFull), "watch-globs", "W", nil, "filepath globs to watch for changes and trigger full regen")
	GenCmd.Flags().StringSliceVarP(&(flags.GenFlags.WatchFast), "watch-fast", "X", nil, "filepath globs to watch for changes and trigger fast regen")
	GenCmd.Flags().StringVarP(&(flags.GenFlags.AsModule), "as-module", "", "", "<github.com/username/<name>> like value for the generator module made from the given flags")
	GenCmd.Flags().StringVarP(&(flags.GenFlags.InitModule), "init", "", "", "<name> to bootstrap a new genarator module")
	GenCmd.Flags().StringVarP(&(flags.GenFlags.Outdir), "outdir", "O", "", "base directory to write all output u")
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

	Short: "modular and composable code gen: CUE & data + templates = _",

	Long: genLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		extra := " "
		// check mode
		if len(flags.GenFlags.Template) == 0 {
			extra = extra + "module"
		} else {
			extra = extra + "adhoc"
		}

		// check watch
		if flags.GenFlags.Watch || len(flags.GenFlags.WatchFull) > 0 || len(flags.GenFlags.WatchFull) > 0 {
			extra = extra + " watch"
		}

		ga.SendCommandPath(cmd.CommandPath() + extra)

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

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	GenCmd.SetHelpFunc(thelp)
	GenCmd.SetUsageFunc(tusage)

}
