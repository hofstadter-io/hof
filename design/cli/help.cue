package cli

import (
	"github.com/hofstadter-io/hof/design/cli/cmds"
)

// TBD:   "Ø"
// TBD:   "α"
// TBD:   "β"
// TBD:   "✓"

#RootCustomHelp: """
hof - a polyglot tool for building software

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]


Initialize and create new hof workspaces:
  \(cmds.#InitCommand.Help)
  \(cmds.#CloneCommand.Help)

Model your designs, generate implementation, run or test anything:
  \(cmds.#DatamodelCommand.Help)
  \(cmds.#GenCommand.Help)
  \(cmds.#RunCommand.Help)
  \(cmds.#RuntimesCommand.Help)
  \(cmds.#TestCommand.Help)

Labels are used _ for _ (see also 'hof topic labels'):
  \(cmds.#LabelCommand.Help)
  \(cmds.#LabelsetCommand.Help)

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
  \(cmds.#CreateCommand.Help)
  \(cmds.#GetCommand.Help)
  \(cmds.#SetCommand.Help)
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

Collaborate (see also 'hof topic collaborate'):
  \(cmds.#FetchCommand.Help)
  \(cmds.#PullCommand.Help)
  \(cmds.#PushCommand.Help)
  \(cmds.#ProposeCommand.Help)
  \(cmds.#PublishCommand.Help)
  \(cmds.#RemotesCommand.Help)

Local development commands:
  \(cmds.#ReproCommand.Help)
  \(cmds.#JumpCommand.Help)
  \(cmds.#UiCommand.Help)
  \(cmds.#TuiCommand.Help)
  \(cmds.#ReplCommand.Help)
  pprof                 go pprof by setting HOF_CPU_PROFILE="hof-cpu.prof" hof <cmd>


Send us feedback or say hello:
  \(cmds.#FeedbackCommand.Help)
                        you can also chat with us on https://gitter.im/hofstadter-io

Additional commands:
  help                  help about any command
  topic                 additional information for various subjects and concepts
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal

Additional topics:
  schema, codegen, modeling, mirgrations
  resources, labels, context, querying
  workflow, changesets, collaboration

(✓) command is generally available
(β) command is beta and ready for testing
(α) command is alpha and under developmenr
(Ø) command is null and yet to be implemented

Flags:
<<flag-usage>>
Use "hof [command] --help / -h" for more information about a command.
Use "hof topic [subject]"  for more information about a subject.

"""
