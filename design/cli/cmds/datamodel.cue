package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#DatamodelCommand: schema.#Command & {
	Name:  "datamodel"
	Usage: "datamodel"
	Aliases: ["dm"]
	Short: "create, view, diff, calculate / migrate, and manage your data models"
	Long:  #DatamodelRootHelp

	OmitRun: true

	Pflags: [...schema.#Flag] & [ {
		Name:    "Datamodels"
		Long:    "datamodel"
		Short:   "d"
		Type:    "[]string"
		Default: "nil"
		Help:    "Datamodels for the datamodel commands"
	}, {
		Name:    "Models"
		Long:    "model"
		Short:   "m"
		Type:    "[]string"
		Default: "nil"
		Help:    "Models for the datamodel commands"
	}, {
		Name:    "Output"
		Long:    "output"
		Short:   "o"
		Type:    "string"
		Default: "\"table\""
		Help:    "Output format [table,cue]"
	}, {
		Name:    "Format"
		Long:    "format"
		Short:   "f"
		Type:    "string"
		Default: "\"_\""
		Help:    "Pick format from Cuetils"
	}, {
		Name:    "Since"
		Long:    "since"
		Short:   "s"
		Type:    "string"
		Default: ""
		Help:    "Timestamp to filter since"
	}, {
		Name:    "Until"
		Long:    "until"
		Short:   "u"
		Type:    "string"
		Default: ""
		Help:    "Timestamp to filter until"
	}]

	Commands: [{
		Name:  "checkpoint"
		Usage: "checkpoint"
		Aliases: ["cp", "x"]
		Short: "create a snapshot of the data model"
		Long:  Short
		Flags: [...schema.#Flag] & [{
			Name:    "bump"
			Long:    "bump"
			Short:   "b"
			Type:    "string"
			Default: "\"patch\""
			Help:    "type of version bump in [major,minor,patch,<semver>]"
		}]
	}, {
		Name:  "diff"
		Usage: "diff"
		Aliases: ["d"]
		Short: "show the diff between data model version"
		Long:  Short
	}, {
		Name:  "history"
		Usage: "history"
		Aliases: ["hist", "h"]
		Short: "list the snapshots for a data model"
		Long:  Short
	}, {
		Name:  "info"
		Usage: "info"
		Aliases: ["i"]
		Short: "print details for a data model"
		Long:  Short
	}, {
		Name:  "list"
		Usage: "list"
		Aliases: ["ls"]
		Short: "print available data models"
		Long:  Short
	}, {
		Name:  "log"
		Usage: "log"
		Aliases: ["l"]
		Short: "show the history of diffs for a data model"
		Long:  Short
	}]
}

#DatamodelRootHelp: """
	Data models are sets of models which are used in many hof processes and modules.
	
	At their core, they represent the most abstract representation for objects and
	their relations in your applications. They are extended and annotated to add
	context fot their usage in different code generators: (DB vs Server vs Client).
	
	Beyond representing models in their current form, a history is maintained so that:
	  - database migrations can be created and managed
	  - servers can handle multiple model versions
	  - clients can implement feature flags
	Much of this is actually handled by code generators and must be implemented there.
	Hof handles the core data model definitions, history, and snapshot creation.
	"""
