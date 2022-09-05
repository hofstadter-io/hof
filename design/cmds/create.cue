package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#CreateCommand: schema.#Command & {
	Name:  "create"
	Usage: "create <module location>"
	Short: "easily bootstrap full project, components, and more"
	Long:  #CreateRootHelp

	Args: [{
		Name: "module"
		Type: "string"
		Required: false
		Help: "git repository or directory with a creator, accepts subdirs on both"
	}]

	Flags: [...schema.#Flag] & [ {
		Name:    "input"
		Long:    "input"
		Short:   "I"
		Type:    "[]string"
		Default: "nil"
		Help:    "inputs to the create module"
	},
	{
		Name:    "generator"
		Type:    "[]string"
		Default: "nil"
		Help:    "generator tags to run, default is all"
		Long:    "generator"
		Short:   "G"
	},
	{
		Name:    "Outdir"
		Type:    "string"
		Default: "\"\""
		Help:    "base directory to write all output to"
		Long:    "outdir"
		Short:   "O"
	}]
}

#CreateRootHelp: #"""
hof create enables you to easily bootstrap
code for full projects, components, and more.

Any generator can support the create command
and most will bootstrap a generator.
This means you get all the same benefits from
hof's code generation engine, turning your
bootstrapped code into a living template.

# create from any git repo and any ref
hof create github.com/username/repo@v1.2.3
hof create github.com/username/repo@a1b2c3f
hof create github.com/username/repo@latest

# -I supplies inputs as key/value pairs or from a file
# when no flag is supplied, an interactive prompt is used
hof create github.com/username/repo@v1.2.3 \
  -I name=foo -I val=bar \
  -I @inputs.cue

# you can also reference local generators by their cue inputs
# the location should start with a '.' (./ or ../) to indicate local mode
hof create ../my-gen
"""#
