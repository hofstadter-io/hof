package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"

	"github.com/hofstadter-io/hof/design/cmds"
)

#Outdir: "./cmd/hof"

#CLI: schema.#Cli & {
	Name:    "hof"
	Package: "github.com/hofstadter-io/hof/cmd/hof"

	Usage:      "hof"
	Short:      "The High Code Framework"
	Long:       Short
	CustomHelp: #RootCustomHelp

	OmitRun:           true
	PersistentPrerun:  true
	PersistentPostrun: true

	Pflags: #CliPflags

	//Commands: [...schema.#Command]
	Commands: [

		// main commands
		cmds.#CreateCommand,
		cmds.#DatamodelCommand,
		cmds.#GenCommand,
		cmds.#FlowCommand,
		cmds.#FmtCommand,
		cmds.#ModCommand,

		// beta commands
		cmds.#RunCommand,

		// additional commands
		cmds.#FeedbackCommand,
	]

	//
	// Addons
	//
	Releases: #CliReleases
	Updates:  true
	EnablePProf: true
	Telemetry: "UA-103579574-5"
}

#RootCustomHelp: """
hof - the high code framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  \(cmds.#CreateCommand.Help)
  \(cmds.#DatamodelCommand.Help)
  \(cmds.#GenCommand.Help)
  \(cmds.#FlowCommand.Help)
  \(cmds.#FmtCommand.Help)
  \(cmds.#ModCommand.Help)

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  \(cmds.#FeedbackCommand.Help)

Flags:
<<flag-usage>>
Use "hof [command] --help / -h" for more information about a command.

"""
