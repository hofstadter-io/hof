package cli

import (
	"github.com/hofstadter-io/hof/design/cli/cmds"
)

// TBD:   "Ø"
// TBD:   "α"
// TBD:   "β"
// TBD:   "✓"
#RootCustomHelp: """
hof - the high code framework

  Learn more at https://docs.hofstadter.io

Usage:
  hof [flags] [command] [args]

Main commands:
  \(cmds.#InitCommand.Help)
  \(cmds.#GenCommand.Help)
  \(cmds.#ModCommand.Help)
  \(cmds.#TestCommand.Help)
  \(cmds.#ConfigCommand.Help)

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

// #OldRootCustomHelp: """
//hof - the high code framework

  //Learn more at https://docs.hofstadter.io

//Usage:
  //hof [flags] [command] [args]


//Initialize and create new hof workspaces:
  //\(cmds.#InitCommand.Help)
  //\(cmds.#CloneCommand.Help)

//Model your designs, generate implementation, run or test anything:
  //\(cmds.#DatamodelCommand.Help)
  //\(cmds.#GenCommand.Help)
  //\(cmds.#RunCommand.Help)
  //\(cmds.#TestCommand.Help)

//Labels are used _ for _ (see also 'hof topic labels'):
  //\(cmds.#LabelCommand.Help)
  //\(cmds.#LabelsetCommand.Help)

//Learn more about hof and the _ you can do:
  //each command has four flags, use 'list' as their arg
	//to see available items on a command
		//--help              print help message
		//--topics            addtional help topics
		//--examples          examples for the command
		//--tutorials         tutorials for the command

//Download modules, add instances or content, and manage runtimes:
  //\(cmds.#ModCommand.Help)
  //\(cmds.#AddCommand.Help)
  //\(cmds.#RuntimesCommand.Help)

//Manage resources (see also 'hof topic resources'):
  //\(cmds.#InfoCommand.Help)
  //\(cmds.#CreateCommand.Help)
  //\(cmds.#GetCommand.Help)
  //\(cmds.#SetCommand.Help)
  //\(cmds.#EditCommand.Help)
  //\(cmds.#DeleteCommand.Help)

//Configure, Unify, Execute (see also https://cuelang.org):
  //\(cmds.#CmdCommand.Help)
  //\(cmds.#DefCommand.Help)
  //\(cmds.#EvalCommand.Help)
  //\(cmds.#ExportCommand.Help)
  //\(cmds.#FormatCommand.Help)
  //\(cmds.#ImportCommand.Help)
  //\(cmds.#TrimCommand.Help)
  //\(cmds.#VetCommand.Help)
  //\(cmds.#StCommand.Help)

//Manage logins, config, secrets, and context:
  //\(cmds.#AuthCommand.Help)
  //\(cmds.#ConfigCommand.Help)
  //\(cmds.#SecretCommand.Help)
  //\(cmds.#ContextCommand.Help)

//Examine workpsace history and state:
  //\(cmds.#StatusCommand.Help)
  //\(cmds.#LogCommand.Help)
  //\(cmds.#DiffCommand.Help)
  //\(cmds.#BisectCommand.Help)

//Grow, mark, and tweak your shared history (see also 'hof topic changesets'):
  //\(cmds.#IncludeCommand.Help)
  //\(cmds.#BranchCommand.Help)
  //\(cmds.#CheckoutCommand.Help)
  //\(cmds.#CommitCommand.Help)
  //\(cmds.#MergeCommand.Help)
  //\(cmds.#RebaseCommand.Help)
  //\(cmds.#ResetCommand.Help)
  //\(cmds.#TagCommand.Help)

//Collaborate (see also 'hof topic collaborate'):
  //\(cmds.#FetchCommand.Help)
  //\(cmds.#PullCommand.Help)
  //\(cmds.#PushCommand.Help)
  //\(cmds.#ProposeCommand.Help)
  //\(cmds.#PublishCommand.Help)
  //\(cmds.#RemotesCommand.Help)

//Local development commands:
  //\(cmds.#ReproCommand.Help)
  //\(cmds.#JumpCommand.Help)
  //\(cmds.#UiCommand.Help)
  //\(cmds.#TuiCommand.Help)
  //\(cmds.#ReplCommand.Help)
  //pprof                 go pprof by setting HOF_CPU_PROFILE="hof-cpu.prof" hof <cmd>


//Send us feedback or say hello:
  //\(cmds.#FeedbackCommand.Help)
                        //you can also chat with us on https://gitter.im/hofstadter-io

//Additional commands:
  //help                  help about any command
  //update                check for new versions and run self-updates
  //version               print detailed version information
  //completion            generate completion helpers for your terminal

//(✓) command is generally available
//(β) command is beta and ready for testing
//(α) command is alpha and under developmenr
//(Ø) command is null and yet to be implemented

//Flags:
//<<flag-usage>>
//Use "hof [command] --help / -h" for more information about a command.
//Use "hof topic [subject]"  for more information about a subject.

//"""
