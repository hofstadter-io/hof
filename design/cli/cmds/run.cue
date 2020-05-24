package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// TODO, JAMStack / HofKit
#RunCommand: schema.#Command & {
	TBD:  "α"
	Name:  "run"
	Usage: "run"
	Aliases: ["r"]
	Short: "run polyglot command and scripts seamlessly across runtimes"
	Long:  Short
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
