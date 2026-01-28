package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

_agentLong: """
	build, use, evaluate, and serve agentic systems
	"""

AgentCommand: schema.Command & {
	Name:  "agent"
	Usage: "agent [...target] [% ...cue]"
	Short: "build, chat with, and serve agentic systems"
	Long:  _agentLong

	PersistentPrerun:     true
	PersistentPrerunBody: "err = runtime.EnsureInfra()"

	OmitRun: false
	Imports: [
		{Path: "github.com/hofstadter-io/hof/lib/runtime"},
		{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"},
	]
	Body: "err = libagentcmd.Main(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)"

	Commands: [{
		Name:  "list"
		Usage: "list [...target] [% ...cue]"
		Short: "list agentic components"
		Long:  "list agentic components"
		Imports: [{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
		Body: "err = libagentcmd.List(args, flags.RootPflags, flags.AgentPflags)"
	},{
		Name:  "chat"
		Usage: "chat [...target] [% ...cue]"
		Short: "chat with an agent"
		Long:  "chat with an agent"
		Imports: [{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
		Body: "err = libagentcmd.Chat(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)"

		Pflags: [{
			Name:    "Sid"
			Long:    "sid"
			Type:    "string"
			Default: "\"\""
			Help:    "session id to continue from"
		}]

		Commands: [{
			Name:  "info"
			Usage: "info [...target] [% ...cue]"
			Short: "get info for session"
			Long:  "get info for session"
			Imports: [{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
			Body: "err = libagentcmd.ChatInfo(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)"
		},{
			Name:  "list"
			Usage: "list [...target] [% ...cue]"
			Short: "list sessions for given targets"
			Long:  "list sessions for given targets"
			Imports: [{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
			Body: "err = libagentcmd.ChatList(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)"
		},{
			Name:  "delete"
			Usage: "delete [...target] [% ...cue]"
			Short: "delete sessions for given targets"
			Long:  "delete sessions for given targets"
			Imports: [{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
			Body: "err = libagentcmd.ChatList(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)"
		}]
	},{
		Name:  "serve"
		Usage: "serve [...target] [% ...cue]"
		Short: "serve agentic components"
		Long:  "serve agentic components"
		Imports: [{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
		Body: "err = libagentcmd.Serve(args, flags.RootPflags, flags.AgentPflags)"
	},{
		Name:  "bulk"
		Usage: "bulk [...target] [% ...cue]"
		Short: "bulk process with agents"
		Long:  "bulk process with agents"
		Imports: [{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
		Body: "err = libagentcmd.Bulk(args, flags.RootPflags, flags.AgentPflags)"
	},{
		Name:  "eval"
		Usage: "eval [...target] [% ...cue]"
		Short: "eval agents over parameters"
		Long:  "eval agents over parameters"
		Imports: [{As: "libagentcmd", Path: "github.com/hofstadter-io/hof/lib/agent/cmd"}, {Path: "github.com/hofstadter-io/hof/cmd/hof/flags"}]
		Body: "err = libagentcmd.Eval(args, flags.RootPflags, flags.AgentPflags)"
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
