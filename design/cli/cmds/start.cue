package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#InitCommand: schema.#Command & {
	TBD:   "β"
	Name:  "init"
	Usage: "init <module>"
	Short: "create an empty workspace or initialize an existing directory to one"
	Long:  """
	create an empty workspace or initialize an existing directory to one

	  module name or path should look like github.com/hofstadter-io/hof
	"""

	Args: [{
		Name:     "module"
		Type:     "string"
		Required: true
		Help:     "module url or path (github.com/hofstadter-io/hof"
	},
	{
		Name:     "name"
		Type:     "string"
		Help:     "module name, defaults to last part of module"
	}]
}

#CloneCommand: schema.#Command & {
	TBD:   "β"
	Name:  "clone"
	Usage: "clone"
	Short: "clone a workspace or repository into a new directory"
	Long:  Short

	Args: [{
		Name:     "module"
		Type:     "string"
		Required: true
		Help:     "module url or path (github.com/hofstadter-io/hof"
	},
	{
		Name:     "name"
		Type:     "string"
		Help:     "module name, defaults to last part of module"
	}]
}
