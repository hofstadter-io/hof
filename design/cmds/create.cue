package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

CreateCommand: schema.Command & {
	Name:  "create"
	Usage: "create <module path> [extra args]"
	Short: "starter kits or blueprints from any git repo"
	Long:  CreateRootHelp

	Args: [{
		Name:     "module"
		Type:     "string"
		Required: true
		Help:     "git repository or directory with a creator, accepts subdirs on both"
	}, {
		Name: "extra"
		Type: "[]string"
		Rest: true
		Help: "extra arguments for the creator, if it accepts them"
	}]

	Flags: [...schema.Flag] & [ {
		Name:    "input"
		Long:    "input"
		Short:   "I"
		Type:    "[]string"
		Default: "nil"
		Help:    "inputs to the create module"
	}, {
		Name:    "generator"
		Type:    "[]string"
		Default: "nil"
		Help:    "generator tags to run, default is all"
		Long:    "generator"
		Short:   "G"
	}, {
		Name:    "Outdir"
		Type:    "string"
		Default: "\"\""
		Help:    "base directory to write all output to"
		Long:    "outdir"
		Short:   "O"
	}, {
		Name:    "exec"
		Type:    "bool"
		Default: "false"
		Help:    "enable pre/post-exec support when generating code"
		Long:    "exec"
	}]
}

CreateRootHelp: #"""
	hof create enables you to easily bootstrap
	code for full projects, components, and more.
	
	Examples can be found in the documentation:
	
	  https://docs.hofstadter.io/hof-create/
	
	By adding one config file and templates to your repo
	your users can quickly bootstrap full applications,
	tooling configuration, and other code using your project.
	Share consistent scaffolding, configurable to users.
	
	Any hof generator can also support the create command
	and most choose to bootstrap a generator at minimum.
	This means you get all the same benefits from
	hof's code generation engine, turning your
	bootstrapped code into a living template.
	
	Run create from any git repo and any ref
	
	  hof create github.com/username/repo@v1.2.3
	  hof create github.com/username/repo@a1b2c3f
	  hof create github.com/username/repo@latest
	
	-I supplies inputs as key/value pairs or from a file
	when no flag is supplied, an interactive prompt is used
	
	  hof create github.com/username/repo@v1.2.3 \
	    -I name=foo -I val=bar \
	    -I @inputs.cue
	
	You can also reference local generators by their cue inputs.
	This local lookup is indicated by ./ or ../ starting a path.
	Use this mode when developing and testing locally.
	
	  hof create ../my-gen
	"""#
