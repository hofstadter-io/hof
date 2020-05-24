package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Kubernetes inspired commands (maybe some hyper-cloud too)

#RuntimesCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "runtimes"
	Usage: "runtimes"
	Short: "work with runtimes"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "+ "
		Name:  "info"
		Usage: "info"
		Short: "print information about known runtimes"
		Long:  Short
	}, {
		TBD:   "+ "
		Name:  "add"
		Usage: "add"
		Short: "add a runtime to your system or workspace"
		Long:  Short
	}, {
		TBD:   "+ "
		Name:  "get"
		Usage: "get"
		Short: "find and display runtimes"
		Long:  Short
	}, {
		TBD:   "+ "
		Name:  "edit"
		Usage: "edit"
		Short: "edit a runtime configuration"
		Long:  Short
	}, {
		TBD:   "+ "
		Name:  "remove"
		Usage: "remove"
		Short: "remove a runtime"
		Long:  Short
	}]
}
