package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#FlowCommand: schema.#Command & {
	Name:  "flow"
	Aliases: ["f"]
	Usage: "flow [cue files...]"
	Short: "run file(s) through the hof/flow DAG engine"
	Long:  """
  \(Short)

  Use hof/flow to transform data, call APIs, work with DBs,
  read and write files, call any program, handle events,
  and much more.

  Docs: https://docs.hofstadter.io/data-flow

  Example:

    @flow()

    call: {
      @task(api.Call)
      req: { ... }
      resp: string
    }

    print: {
      @task(os.Stdout)
      test: call.resp
    }

  """

	Args: [{
		Name:     "globs"
		Type:     "[]string"
		Help:     "file globs to the operation"
		Rest:			true
	}]

	Flags: [{
		Name:    "list"
		Long:    "list"
		Short:   "l"
		Type:    "bool"
    Default: "false"
		Help:    "list available pipelines"
	},{
		Name:    "docs"
		Long:    "docs"
		Short:   "d"
		Type:    "bool"
    Default: "false"
		Help:    "print pipeline docs"
	},{
		Name:    "flow"
		Long:    "flow"
		Short:   "f"
		Type:    "[]string"
    Default: "nil"
		Help:    "flow labels to match and run"
	},{
		Name:    "tags"
		Long:    "tags"
		Short:   "t"
		Type:    "[]string"
    Default: "nil"
		Help:    "data tags to inject before run"
	},{
		Name:    "DebugTasks"
		Long:    "debug-tasks"
		Short:   ""
		Type:    "bool"
    Default: "false"
		Help:    "print debugging info about tasks"
	}, {
    Name:    "stats"
    Type:    "bool"
    Default: "false"
    Help:    "Print final task statistics"
    Long:    "stats"
    Short:   "s"
  }]
}
