package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"

	"github.com/hofstadter-io/hof/design/cmds"
)

Outdir: "./cmd/hof"

CLI: schema.Cli & {
	Name:   "hof"
	Module: "github.com/hofstadter-io/hof"

	Usage:      "hof"
	Short:      "The Higher Order Framework"
	Long:       Short
	CustomHelp: RootCustomHelp

	OmitRun:           true
	PersistentPrerun:  true
	PersistentPostrun: true

	Pflags: CliPflags

	//Commands: [...schema.Command]
	Commands: [

		// main commands
		cmds.CreateCommand,
		cmds.DatamodelCommand,
		cmds.GenCommand,
		cmds.FlowCommand,
		cmds.FmtCommand,
		cmds.ModCommand,

		// cue commands
		cmds.DefCommand,
		cmds.EvalCommand,
		cmds.ExportCommand,
		cmds.VetCommand,

		// beta commands (tied to the development of custom coding agent for vscode)
		cmds.ExtensionCommand,
		cmds.AgentCommand,
		cmds.EnvCommand,

		// alpha commands
		cmds.RunCommand,
		cmds.TuiCommand,

		// additional commands
		cmds.FeedbackCommand,
	]

	//
	// Addons
	//
	Releases:    CliReleases
	Updates:     true
	EnablePProf: true
	Telemetry:   "G-6CYEVMZL4R"
}

RootCustomHelp: """
hof - the higher order framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  \(cmds.AgentCommand.Help)
  \(cmds.CreateCommand.Help)
  \(cmds.DatamodelCommand.Help)
  \(cmds.DefCommand.Help)
  \(cmds.EnvCommand.Help)
  \(cmds.EvalCommand.Help)
  \(cmds.ExportCommand.Help)
  \(cmds.ExtensionCommand.Help)
  \(cmds.FlowCommand.Help)
  \(cmds.FmtCommand.Help)
  \(cmds.GenCommand.Help)
  \(cmds.ModCommand.Help)
  \(cmds.TuiCommand.Help)
  \(cmds.VetCommand.Help)

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  \(cmds.FeedbackCommand.Help)

Flags:
<<flag-usage>>
Use "hof [command] --help / -h" for more information about a command.

"""
