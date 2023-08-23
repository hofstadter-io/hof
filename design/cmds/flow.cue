package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

FlowCommand: schema.Command & {
	Name: "flow"
	Aliases: ["f"]
	Usage: "flow [cue files...] [@flow/name...] [+key=value]"
	Short: "run workflows and tasks powered by CUE"
	Long:  """
  \(Short)

  Use hof/flow to transform data, call APIs, work with DBs,
  read and write files, call any program, handle events,
  and much more.

  'hof flow' is very similar to 'cue cmd' and built on the same flow engine.
  Tasks and dependencies are inferred.
  Hof flow has a slightly different interface and more task types.

  Docs: https://docs.hofstadter.io/data-flow

  Example:

    @flow()

    call: {
      @task(api.Call)
      req: { ... }
      resp: {
        statusCode: 200
        body: string
      }
    }

    print: {
      @task(os.Stdout)
      test: call.resp
    }

  Arguments:
    cue entrypoints are the same as the cue cli
    @path/name  is shorthand for -f / --flow should match the @flow(path/name)
    +key=value  is shorthand for -t / --tags and are the same as CUE injection tags

    arguments can be in any order and mixed

  @flow() indicates a flow entrypoint
    you can have many in a file or nested values
    you can run one or many with the -f flag

  @task() represents a unit of work in the flow dag
    intertask dependencies are autodetected and run appropriately
    hof/flow provides many built in task types
    you can reuse, combine, and share as CUE modules, packages, and values

  """

	Args: [{
		Name: "globs"
		Type: "[]string"
		Help: "file globs to the operation"
		Rest: true
	}]

	Flags: [{
		Name:    "list"
		Long:    "list"
		Type:    "bool"
		Default: "false"
		Help:    "list available pipelines"
	}, {
		Name:    "flow"
		Long:    "flow"
		Short:   "f"
		Type:    "[]string"
		Default: "nil"
		Help:    "flow labels to match and run"
	}, {
		Name:    "progress"
		Long:    "progress"
		Short:   ""
		Type:    "bool"
		Default: "false"
		Help:    "print task progress as it happens"
	}]
}
