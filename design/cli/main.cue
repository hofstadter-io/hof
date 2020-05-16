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

	OmitRun: true

	Imports: [
		{Path: "github.com/hofstadter-io/hof/lib/runtime"},
	]

	PersistentPrerun: true
	PersistentPrerunBody: "runtime.Init()"

  PersistentPostrun: true

  Pflags: #CliPflags

  Commands: [...schema.#Command] & [
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
		cmds.#DevCommand,
		cmds.#UiCommand,
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

  // EnablePProf: true

}

