package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

_envLong: """
	build, run, ship, and deploy environments (image, service, stack)
	
	'veg env' looks for custom commands and treats them equally to builtin commands.
	All commands have well known behaviors, depending on the $kind of a target.
	'veg env list' and 'hof env sync' work with all '$kind's.
	Most commands only work with a subset that makes sense or is explicit.
	See their help text to learn more.
	
	Note, hof -> veg in terms of name, I have started with this command, and also the agentic stuff
	The big move is coming soon...
	
	## Examples
	
	# [...targets], or "points", are selected via args and flags
	#   list allows you to explore that space without syncing or triggering evaluation
	veg env list|info ['^name$'] [-K '^kind$'] [-P '^path$'] [-S name|kind|path]
	
	# sync, evaluates the DAGs, but doesn't export or server anything, steps are run
	#   it can act as "real dry run" compared to the --dry-run flag for other commands
	#   without any args or flags, sync acts as a big test as well
	veg env sync [...targets] [...flags]
	
	# run, creates an interactive session and binds and dependent services
	#   this is closest to docker run or kubectl exec
	veg env run [...target] [...flags]
	
	# run, launches an services or stacks, similar to compose and helm
	veg env up [...target] [...flags]
	
	# export artifacts, local or remote, object storage and registries
	veg env export -P release -T v0.4.3 -t dest=./release
	
	# make your own commands and flags, designed for your workflows
	#   this is closest to Makefiles or package.json scripts
	#   define similar commands with the power of CUE and Dagger
	veg env [init, test, lint, ci, publish, deploy, ...]
	veg env ... -t env=stg -t stack=app -t branch=main
	
	## Important References
	
	./schemas/env    # the CUE schemas for what you can do in veg/env
	./catalogs/env   # reusable CUE for all sorts of things from small to big
	./examples/env   # simple and complex examples for you to play and fork
	
	"""

EnvCommand: schema.Command & {
	Name:  "env"
	Usage: "env [...target] [% ...cue]"
	Short: "build, run, ship, and deploy environments (image, service, stack)"
	Long:  _envLong

	Imports: [
		{Path: "github.com/hofstadter-io/hof/lib/runtime"},
		{As: "libenvcmd", Path: "github.com/hofstadter-io/hof/lib/env/cmd"},
	]

	PersistentPrerun:     true
	PersistentPrerunBody: "err = runtime.EnsureInfra()"

	// this runs the commands or naked `veg env` command
	Body: "err = libenvcmd.Env(args, flags.RootPflags, flags.EnvPflags)"

	Commands: [{
		Name:  "sync"
		Usage: "sync [...target] [% ...cue]"
		Short: "sync target points in an environment"
		Long:  "sync target points in an environment, making sure they are ready to go, no matter the type"
		Imports: [{As: "libenvcmd", Path: "github.com/hofstadter-io/hof/lib/env/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
		Body: "err = libenvcmd.Sync(args, flags.RootPflags, flags.EnvPflags)"
	}, {
		Name:  "export"
		Usage: "export [...target] [% ...cue]"
		Short: "export target points from an environment to outside world"
		Long:  "export target points from an environment to outside world, for each point .. for each tag"
		Imports: [{As: "libenvcmd", Path: "github.com/hofstadter-io/hof/lib/env/cmd"}]
		Body: "err = libenvcmd.Export(args, flags.RootPflags, flags.EnvPflags, flags.Env__ExportFlags)"
		Flags: [{
			Name:    "Tag"
			Long:    "tag"
			Short:   "T"
			Type:    "[]string"
			Default: "nil" // todo, support special options like git-tag or git-commit "auto" that has an understanding of where it is running (list out the handful of variables that differentiate between env's env (local, ci, deployed), which each can have any user defined params as well)
			Help:    "tags to give to the environment, can be set multiple times"
		}]
	}, {
		Name:  "info"
		Usage: "info [...target] [% ...cue]"
		Short: "get details for target points in an environments"
		Long:  "get details for target points in an environments"
	}, {
		Name:  "list"
		Usage: "list [...target] [% ...cue]"
		Short: "list points in an environment"
		Long:  "list points in an environment"
		Imports: [{As: "libenvcmd", Path: "github.com/hofstadter-io/hof/lib/env/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
		Body: "err = libenvcmd.List(args, flags.RootPflags, flags.EnvPflags)"
	}, {
		Name:  "run"
		Usage: "run <target> [% [...cue]]"
		Short: "run target point in an environment"
		Long:  "run target point in an environment"
		Imports: [{As: "libenvcmd", Path: "github.com/hofstadter-io/hof/lib/env/cmd"}]
		Body: "err = libenvcmd.Run(args, flags.RootPflags, flags.EnvPflags, flags.Env__RunFlags)"
		Flags: [{
			Name:    "Command"
			Long:    "cmd"
			Short:   "c"
			Type:    "string"
			Default: "\"\""
			Help:    "the command to run, if none by default or to override"
		}]
	}, {
		Name:  "up"
		Usage: "up [...target] [% ...cue]"
		Short: "starts target points in an environment"
		Long:  "starts target points in an environment, this is very similar to docker-compose or helm locally"
		Imports: [{As: "libenvcmd", Path: "github.com/hofstadter-io/hof/lib/env/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
		Body: "err = libenvcmd.Up(args, flags.RootPflags, flags.EnvPflags)"
	}]

	Pflags: [...schema.Flag] & [{
		Name:    "Renderer"
		Long:    "renderer"
		Short:   "R"
		Type:    "string"
		Default: "\"auto\""
		Help:    "output format [auto, plain, tty, dots, report (for ai)]"
	}, {
		Name:    "OnFailure"
		Long:    "on-failure"
		Short:   "F"
		Type:    "bool"
		Default: "false"
		Help:    "on failure, enter an interactive terminal, requires a tty"
	}, {
		Name:    "NoExit"
		Long:    "no-exit"
		Short:   "N"
		Type:    "bool"
		Default: "false"
		Help:    "Leave the TUI open after finishing"
	}, {
		Name:    "NoCache"
		Long:    "no-cache"
		Short:   "Z"
		Type:    "bool"
		Default: "false"
		Help:    "bust the cache and force evaluation"
	}, {
		Name:    "Kind"
		Long:    "kind"
		Short:   "K"
		Type:    "[]string"
		Default: "nil"
		Help:    "kinds to include, defaults to all"
	}, {
		Name:    "Sort"
		Long:    "sort"
		Short:   "S"
		Type:    "[]string"
		Default: "nil"                       // todo, support special options like git-tag or git-commit
		Help:    "sort columns, can be used multiple times" // todo, support +/- prefix for asc/desc
	}, {
		Name:    "EnvVar"
		Long:    "env-var"
		Type:    "[]string"
		Default: "nil"
		Help:    "key=value ENV vars to pass"
	}, {
		Name:    "EnvFile"
		Long:    "env-file"
		Type:    "[]string"
		Default: "nil"
		Help:    "path to a file with ENV vars to pass"
	}, {
		Name:    "ShhVar"
		Long:    "shh-var"
		Type:    "[]string"
		Default: "nil"
		Help:    "key=value secret ENV vars to pass"
	}, {
		Name:    "ShhFile"
		Long:    "shh-file"
		Type:    "[]string"
		Default: "nil"
		Help:    "path to a file with secret ENV vars to pass"
	}, {
		Name:    "EnvAll"
		Long:    "env-all"
		Type:    "bool"
		Default: "false"
		Help:    "pass os.Env (everything)"
	}, {
		Name:    "ShowAll"
		Long:    "all"
		Type:    "bool"
		Default: "false"
		Help:    "show all env targets, not just @env() ones"
	}, {
		Name:    "FailFast"
		Long:    "fail-fast"
		Type:    "bool"
		Default: "false"
		Help:    "fail at first error instead of attempting all targets"
	}, {
		Name:    "Unsafe"
		Long:    "unsafe"
		Type:    "bool"
		Default: "false"
		Help:    "set insecure root capabilities and privileged nesting, use at your own risk, needed for inception"
	}, {
		Name:    "parallel"
		Long:    "parallel"
		Short:   "P"
		Type:    "int"
		Default: "1"
		Help:    "number of args or objects to process at once, they may be highly parallel internally"
	}]

}
