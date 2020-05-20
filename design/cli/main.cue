package cli

import (
	"github.com/hofstadter-io/hofmod-cli/schema"

	"github.com/hofstadter-io/hof/design/cli/cmds"
)

#Outdir: "./cmd/hof"

#CLI: schema.#Cli & {
	Name:    "hof"
	Package: "github.com/hofstadter-io/hof/cmd/hof"

	Usage:      "hof"
	Short:      "Polyglot Code Gereration Framework"
	Long:       Short
	CustomHelp: #RootCustomHelp

	OmitRun: true

	Imports: [
		{Path: "github.com/hofstadter-io/hof/lib/runtime"},
	]

	PersistentPrerun:     true
	PersistentPrerunBody: "runtime.Init()"

	PersistentPostrun: true

	Pflags: #CliPflags

	Commands: [
		// start
		cmds.#InitCommand,
		cmds.#CloneCommand,

		// hof
		cmds.#GenCommand,
		cmds.#ModelsetCommand,
		// #EtlCommand,

		// hof + cue
		cmds.#ModCommand,
		cmds.#AddCommand,
		cmds.#CmdCommand,

		// resources
		cmds.#LabelCommand,
		cmds.#CreateCommand,
		cmds.#ApplyCommand,
		cmds.#GetCommand,
		cmds.#DeleteCommand,

		// cue
		cmds.#DefCommand,
		cmds.#EvalCommand,
		cmds.#ExportCommand,
		cmds.#FormatCommand,
		cmds.#ImportCommand,
		cmds.#TrimCommand,
		cmds.#VetCommand,
		cmds.#StCommand,

		// base
		cmds.#AuthCommand,
		cmds.#ConfigCommand,
		cmds.#SecretCommand,
		cmds.#ContextCommand,

		// workspace / workflow / git commands
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
		// publish
		// release

		// additional help topics
		cmds.#TopicCommand,
		cmds.#FeedbackCommand,
		// bugreport
		// crashreport
		// changelog --version

		// dev & more st commands
		cmds.#DocCommand,
		cmds.#JumpCommand,
		cmds.#BuildCommand,
		// lint
		// fmt
		// fix
		// simplify
		// test
		// bench
		// scan
		// note / knowledge graph
		// todo / scrum


		// TODO: SECURITY
		// - report
		// - scan
		// - fix

		cmds.#UiCommand,
		cmds.#TuiCommand,
		cmds.#ReplCommand,

		// hacks down this way
		{
			Hidden: true
			Name:   "hack"
			Usage:  "hack ..."
			Short:  "development command"
			Long:   Short
		},
		cmds.#GebCommand,
		cmds.#LogoCommand,
	]

	//
	// Addons
	//
	Releases: #CliReleases
	Updates:  true

	Telemetry:      "UA-103579574-5"
	TelemetryIdDir: "hof"

	EnablePProf: true
}
