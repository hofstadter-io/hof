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
	Short:      "The High Code Framework"
	Long:       Short
	CustomHelp: #RootCustomHelp

	OmitRun: true
	PersistentPrerun:     true
	PersistentPostrun: true

	Pflags: #CliPflags

	//Commands: [...schema.#Command]
	Commands: [

		// main commands
		cmds.#GenCommand,
		cmds.#TestCommand,
		cmds.#ModCommand,

		// beta commands
		cmds.#DatamodelCommand,
		cmds.#RunCommand,

		// additional commands
		cmds.#FeedbackCommand,
		// hacks down this way
		{
			Hidden: true
			Name:   "hack"
			Usage:  "hack ..."
			Aliases:  ["h", "x"]
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

	EnablePProf: true
}

#RootCustomHelp: """
hof - the high code framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  \(cmds.#GenCommand.Help)
  \(cmds.#ModCommand.Help)
  \(cmds.#TestCommand.Help)

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  \(cmds.#FeedbackCommand.Help)

Flags:
<<flag-usage>>
Use "hof [command] --help / -h" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.

"""
