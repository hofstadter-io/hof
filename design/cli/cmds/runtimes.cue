package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Manage runtimes for many popular systems and languages.
#RuntimesCommand: schema.#Command & {
	TBD:   "α"
	Name:  "runtimes"
	Usage: "runtimes"
	Short: "work with runtimes (go, js, py, bash, docker, cloud-vms, k8s, custom)"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "α"
		Name:  "info"
		Usage: "info"
		Aliases: ["i"]
		Short: "print information about known runtimes"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "create"
		Usage: "create"
		Aliases: ["c"]
		Short: "add a runtime to your system or workspace"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "get"
		Usage: "get"
		Aliases: ["g"]
		Short: "find and display runtime configurations"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "set"
		Usage: "set"
		Aliases: ["s"]
		Short: "find and configure runtimes"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "edit"
		Usage: "edit"
		Aliases: ["e"]
		Short: "edit a runtime configuration"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "delete"
		Usage: "delete"
		Aliases: ["del"]
		Short: "delete a runtime configuration"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "install"
		Usage: "install"
		Aliases: ["I"]
		Short: "install a runtime"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "uninstall"
		Usage: "uninstall"
		Short: "uninstall a runtime"
		Long:  Short
	}]
}

//#RuntimesCommand: schema.#Command & {
//Name:  "runtimes"
//Usage: "runtimes"
//Aliases: ["r"]
//Short: "manage installed runtimes, versions, and contexts"
//Long: """
//manage installed runtimes, versions, and contexts
//"""

//Flags: [...schema.#Flag] & [
//{
//Name:    "stats"
//Type:    "bool"
//Default: "false"
//Help:    "Print generator statistics"
//Long:    "stats"
//Short:   ""
//},
//{
//Name:    "generator"
//Type:    "[]string"
//Default: "nil"
//Help:    "Generators to run, default is all discovered"
//Long:    "generator"
//Short:   "g"
//},
//]
//}

//#BuildpacksCommand: schema.#Command & {
//Name:  "buildpacks"
//Usage: "buildpacks"
//Aliases: ["r"]
//Short: "manage installed buildpacks, versions, and contexts"
//Long: """
//generate all the things, from code to data to config...
//"""

//Flags: [...schema.#Flag] & [
//{
//Name:    "stats"
//Type:    "bool"
//Default: "false"
//Help:    "Print generator statistics"
//Long:    "stats"
//Short:   ""
//},
//{
//Name:    "generator"
//Type:    "[]string"
//Default: "nil"
//Help:    "Generators to run, default is all discovered"
//Long:    "generator"
//Short:   "g"
//},
//]
//}
