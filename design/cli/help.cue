package cli

import (
	"github.com/hofstadter-io/hof/design/cli/cmds"
)

#RootCustomHelp: """
hof - a polyglot tool for building software

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]


Setup hof and create workspaces and datasets:
  \(cmds.#SetupCommand.Help)
  \(cmds.#InitCommand.Help)
  \(cmds.#CloneCommand.Help)

Model your designs, generate implementation, run anything:
  \(cmds.#ModelsetCommand.Help)
  \(cmds.#GenCommand.Help)
  \(cmds.#RunCommand.Help)
  \(cmds.#RuntimesCommand.Help)

Learn more about hof and the _ you can do:
  \(cmds.#DocCommand.Help)
  \(cmds.#TourCommand.Help)
  \(cmds.#TutorialCommand.Help)


Download modules, add content, and execute commands:
  \(cmds.#ModCommand.Help)
  \(cmds.#AddCommand.Help)
  \(cmds.#CmdCommand.Help)

Manage resources (see also 'hof topic resources'):
  \(cmds.#InfoCommand.Help)
  \(cmds.#LabelCommand.Help)
  \(cmds.#CreateCommand.Help)
  \(cmds.#ApplyCommand.Help)
  \(cmds.#GetCommand.Help)
  \(cmds.#EditCommand.Help)
  \(cmds.#DeleteCommand.Help)

Configure, Unify, Execute (see also https://cuelang.org):
  (also a whole bunch of other awesome things)
  \(cmds.#DefCommand.Help)
  \(cmds.#EvalCommand.Help)
  \(cmds.#ExportCommand.Help)
  \(cmds.#FormatCommand.Help)
  \(cmds.#ImportCommand.Help)
  \(cmds.#TrimCommand.Help)
  \(cmds.#VetCommand.Help)
  \(cmds.#StCommand.Help)


Manage logins, config, secrets, and context:
  \(cmds.#AuthCommand.Help)
  \(cmds.#ConfigCommand.Help)
  \(cmds.#SecretCommand.Help)
  \(cmds.#ContextCommand.Help)

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
  \(cmds.#ReproCommand.Help)
 
Local development commands:
  \(cmds.#JumpCommand.Help)
  \(cmds.#BuildCommand.Help)
  \(cmds.#ReproCommand.Help)
  \(cmds.#UiCommand.Help)
  \(cmds.#TuiCommand.Help)
  \(cmds.#ReplCommand.Help)
  pprof         Go pprof by setting HOF_CPU_PROFILE="hof-cpu.prof" hof <cmd>


Send us feedback or say hello:
  \(cmds.#FeedbackCommand.Help)

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

(+) command is yet to be implemented

Flags:
<<flag-usage>>
Use "hof [command] --help" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.

"""
