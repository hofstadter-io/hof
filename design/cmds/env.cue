package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

EnvCommand: schema.Command & {
	Name:  "env"
	Usage: "env [args]"
	Short: "build, run, ship, and deploy environments (image, service, stack)"
	Long:  "build, run, ship, and deploy environments (image, service, stack)"
	OmitRun: true

	// Pflags: [...schema.Flag] & [ {
	// 	Name:    "Datamodels"
	// 	Long:    "model"
	// 	Short:   "M"
	// 	Type:    "[]string"
	// 	Default: "nil"
	// 	Help:    "specify one or more data models to operate on"
	// }, {

	// all subcommands get the name arg
	// Commands: [...{
	// 	Args: [{
	// 		Name:     "name"
	// 		Type:     "string"
	// 		Required: true
	// 		Help:     "name of the environment"
	// 	}, ..._]
	// }]

	Commands: [{
		Name:  "build"
		Usage: "build <name>"
		Short: "build an environment"
		Long: "build an environment"
	},{
		Name:  "info"
		Usage: "info <name>"
		Short: "get info about an environments"
		Long: "get info about an environments"
	},{
		Name:  "list"
		Usage: "list <name-pattern>"
		Short: "list environments"
		Long: "list environments"
	},{
		Name:  "ps"
		Usage: "ps <name-pattern>"
		Short: "print stats for running environments"
		Long: "print stats for running environments"
	},{
		Name:  "run"
		Usage: "run <name>"
		Short: "run an environment"
		Long: "run an environment"
	},{
		Name:  "push"
		Usage: "push <name>"
		Short: "push an environment"
		Long: "push an environment"
	},{
		Name:  "pull"
		Usage: "pull <name>"
		Short: "pull an environment"
		Long: "pull an environment"
	},{
		Name:  "deploy"
		Usage: "deploy <name>"
		Short: "deploy an environment"
		Long: "deploy an environment"
	}]
}

