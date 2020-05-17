package cli

import (
	"github.com/hofstadter-io/hofmod-cli/schema"

	"github.com/hofstadter-io/hof/design/cli/cmds"
)

#Outdir: "./cmd/hof"

#CLI: schema.#Cli & {
	Name:    "hof"
	Package: "github.com/hofstadter-io/hof/cmd/hof"

	Usage: "hof"
	Short: "Polyglot Code Gereration Framework"
	Long:  Short
	CustomHelp: #RootCustomHelp

	OmitRun: true

	Imports: [
		{Path: "github.com/hofstadter-io/hof/lib/runtime"},
	]

	PersistentPrerun: true
	PersistentPrerunBody: "runtime.Init()"

  PersistentPostrun: true

  Pflags: #CliPflags

  Commands: [
		// base
		cmds.#AuthCommand,
		cmds.#ConfigCommand,
		cmds.#SecretCommand,

		// workspace / workflow / git commands
		cmds.#CloneCommand,
		cmds.#InitCommand,

		cmds.#StatusCommand,
		cmds.#LogCommand,
		cmds.#DiffCommand,
		cmds.#BisectCommand,

		cmds.#IncludeCommand,
		cmds.#BranchCommand,
		cmds.#CheckoutCommand,
		cmds.#CommitCommand,
		cmds.#MergeCommand,
		cmds.#RebaseCommand,
		cmds.#ResetCommand,
		cmds.#TagCommand,

		cmds.#FetchCommand,
		cmds.#PullCommand,
		cmds.#PushCommand,
		cmds.#ProposeCommand,

		// hof
		cmds.#ModCommand,
		cmds.#GenCommand,
		cmds.#ModelCommand,
		// #EtlCommand,

		// hof + cue
		cmds.#AddCommand,
		cmds.#CmdCommand,

		// cue
		cmds.#DefCommand,
		cmds.#EvalCommand,
		cmds.#ExportCommand,
		cmds.#FormatCommand,
		cmds.#ImportCommand,
		cmds.#TrimCommand,
		cmds.#VetCommand,

		// resources
		cmds.#LabelCommand,
		cmds.#CreateCommand,
		cmds.#ApplyCommand,
		cmds.#GetCommand,
		cmds.#DeleteCommand,

		// additional help topics
		cmds.#TopicCommand,

		// dev & more st commands
		cmds.#UiCommand,
		cmds.#TuiCommand,
		cmds.#ReplCommand,

		{
			Hidden: true
			Name:    "hack"
			Usage:   "hack ..."
			Short:   "development command"
			Long: Short
		},
	]

	//
	// Addons
	//
	Releases: #CliReleases
	Updates: true

  Telemetry: "UA-103579574-5"
  TelemetryIdDir: "hof"

  EnablePProf: true

}

// Long Strings

#RootCustomHelp: """
hof - a polyglot tool for building software

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]


Create or clone workspaces and datasets:
  \(cmds.#CloneCommand.Help)
  \(cmds.#InitCommand.Help)

Model your world and generate implementation:
  \(cmds.#ModelCommand.Help)
  \(cmds.#GenCommand.Help)

Download modules, add content, and run commands:
  \(cmds.#ModCommand.Help)
  \(cmds.#AddCommand.Help)
  \(cmds.#CmdCommand.Help)

Configure, Unify, Execute (see also https://cuelang.org):
  \(cmds.#DefCommand.Help)
  \(cmds.#EvalCommand.Help)
  \(cmds.#ExportCommand.Help)
  \(cmds.#FormatCommand.Help)
  \(cmds.#ImportCommand.Help)
  \(cmds.#TrimCommand.Help)
  \(cmds.#VetCommand.Help)

Manage resources (see also 'hof topic resources'):
  \(cmds.#LabelCommand.Help)
  \(cmds.#CreateCommand.Help)
  \(cmds.#ApplyCommand.Help)
  \(cmds.#GetCommand.Help)
  \(cmds.#DeleteCommand.Help)

Manage logins, config, and secrets:
  \(cmds.#AuthCommand.Help)
  \(cmds.#ConfigCommand.Help)
  \(cmds.#SecretCommand.Help)

Examine workpsace history and state:
  \(cmds.#StatusCommand.Help)
  \(cmds.#LogCommand.Help)
  \(cmds.#DiffCommand.Help)
  \(cmds.#BisectCommand.Help)

Grow, mark, and tweak your shared history (see also 'hof topic changesets'):
  \(cmds.#IncludeCommand.Help)
  \(cmds.#BranchCommand.Help)
  \(cmds.#CheckoutCommand.Help)
  \(cmds.#CommitCommand.Help)
  \(cmds.#MergeCommand.Help)
  \(cmds.#RebaseCommand.Help)
  \(cmds.#ResetCommand.Help)
  \(cmds.#TagCommand.Help)

Colloaborate (see also 'hof topic collaborate'):
  \(cmds.#FetchCommand.Help)
  \(cmds.#PullCommand.Help)
  \(cmds.#PushCommand.Help)
  \(cmds.#ProposeCommand.Help)
 
Local development commands:
  \(cmds.#UiCommand.Help)
  \(cmds.#TuiCommand.Help)
  \(cmds.#ReplCommand.Help)
  (run pprof)   HOF_CPU_PROFILE="hof-cpu.prof" hof ...

Additional commands:
  help          Help about any command
  topic         Additional information for various subjects and concepts
  update        Check for new versions and run self-updates
  version       Print detailed version information
  completion    Generate completion helpers for your terminal

Additional topics:
  schema, codegen, modeling, mirgrations
  resources, labels, context, querying
  workflow, changesets, collaboration

Flags:
<<flag-usage>>
Use "hof [command] --help" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.

"""
