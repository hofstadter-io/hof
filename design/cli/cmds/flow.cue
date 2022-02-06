package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#FlowCommand: schema.#Command & {
	Name:  "flow"
	Aliases: ["f"]
	Usage: "flow [cue files...]"
	Short: "run file(s) through the hof/flow DAG engine"
	Long:  Short

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
	}]
}
