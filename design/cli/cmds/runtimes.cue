package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Kubernetes inspired commands (maybe some hyper-cloud too)

#RuntimesCommand: schema.#Command & {
	TBD:   "α"
	Name:  "runtimes"
	Usage: "runtimes"
	Short: "work with runtimes (go, js, py, bash, custom)"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "α"
		Name:  "info"
		Usage: "info"
		Short: "print information about known runtimes"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "add"
		Usage: "add"
		Short: "add a runtime to your system or workspace"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "get"
		Usage: "get"
		Short: "find and display runtimes"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "edit"
		Usage: "edit"
		Short: "edit a runtime configuration"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "remove"
		Usage: "remove"
		Short: "remove a runtime"
		Long:  Short
	}]
}
