package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#DatamodelCommand: schema.#Command & {
	Name:  "datamodel"
	Usage: "datamodel"
	Aliases: ["dm"]
	Short: "manage, diff, and migrate your data models"
	Long:  #DatamodelRootHelp

	OmitRun: true

	Pflags: [...schema.#Flag] & [ {
		Name:    "Datamodels"
		Long:    "datamodel"
		Short:   "d"
		Type:    "[]string"
		Default: "nil"
		Help:    "specify one or more datamodels to operate on"
	}, {
		// TODO, move this up to root pflags
		Name:    "Expression"
		Long:    "expr"
		Short:   "e"
		Type:    "[]string"
		Default: "nil"
		Help:    "CUE paths to select outputs, depending on the command"
		//}, {
		//Name:    "Output"
		//Long:    "output"
		//Short:   "o"
		//Type:    "string"
		//Default: "\"table\""
		//Help:    "Output format [table,cue]"
		//}, {
		//Name:    "Format"
		//Long:    "format"
		//Short:   "f"
		//Type:    "string"
		//Default: "\"_\""
		//Help:    "Pick format from Cuetils"
		//}, {
		//Name:    "After"
		//Long:    "after"
		//Short:   "a"
		//Type:    "string"
		//Default: ""
		//Help:    "Timestamp or version to filter with"
		//}, {
		//Name:    "Before"
		//Long:    "before"
		//Short:   "b"
		//Type:    "string"
		//Default: ""
		//Help:    "Timestamp to filter to filter with"
	}]

	Commands: [{
		Name:  "checkpoint"
		Usage: "checkpoint"
		Aliases: ["cp", "x"]
		Short: "create a snapshot of the data model"
		Long:  Short
		Flags: [...schema.#Flag] & [{
			//Name:    "bump"
			//Long:    "bump"
			//Short:   "B"
			//Type:    "string"
			//Default: "\"patch\""
			//Help:    "type of version bump in [major,minor,patch,<semver>]"
			//}, {
			Name:    "message"
			Long:    "message"
			Short:   "m"
			Type:    "string"
			Default: "\"\""
			Help:    "message describing the checkpoint"
		}]
	}, {
		Name:  "diff"
		Usage: "diff"
		Aliases: ["d"]
		Short: "show the current diff or between datamodel versions"
		Long:  Short
	}, {
		Name:  "tree"
		Usage: "tree"
		Aliases: ["t"]
		Short: "print datamodel structure as a tree"
		Long:  Short
	}, {
		Name:  "list"
		Usage: "list"
		Aliases: ["ls"]
		Short: "print available datamodels"
		Long:  Short
	}, {
		Name:  "log"
		Usage: "log"
		Aliases: ["l"]
		Short: "show the history for a datamodel"
		Long:  Short
		Flags: [...schema.#Flag] & [{
			Name:    "ByValue"
			Long:    "by-value"
			Short:   ""
			Type:    "bool"
			Default: "false"
			Help:    "display snapshot log by value"
		}, {
			Name:    "details"
			Long:    "details"
			Short:   ""
			Type:    "bool"
			Default: "false"
			Help:    "print more when displaying the log"
		}]
	}]
}

#DatamodelRootHelp: """
	Data models are values or objects used in many hof processes and modules.
	The "datamodel" command helps you manage them and track their change history.
	At their core, they represent the models that make up your application.
	The intention is to define a data model for your entire application once,
	then use this source of truth to generate code from database to server to client.
	
	Hof's schema for datamodels is minimal and flexible, allowing you to define the
	shape based on your application. You can have multiple datamodels as well.
	You can also control and where and how history should be tracked. This history
	is included during code generation so that database migrations and functions
	for converting between versions can be created.
	
	# Examples Datamodels

	-- config.cue --
	package datamodel

	import (
		"github.com/hofstadter-io/hof/schema/dm"
		"github.com/hofstadter-io/hof/schema/dm/fields"
	)

	// Track an entire oject
	Config: dm.Object & {

		host: fields.String & { Default: "8080" }

		database: {
			host:   fields.String
			port:   fields.String
			dbconn: fields.String
		}
	}

	-- database.cue --
	package datamodel

	import (
		"github.com/hofstadter-io/hof/schema/dm/sql"
		"github.com/hofstadter-io/hof/schema/dm/fields"
	)

	// Traditional database model which maps onto tables & columns
	Datamodel: sql.Datamodel & {
		// implied through definition, duplicated here for example clarity
		$hof: metadata: {
			id:   "datamodel-abc123"
			name: "MyDatamodel"
		}

		Models: {
			User: {
				Fields: {
					ID:        fields.UUID
					CreatedAt: fields.Datetime
					UpdatedAt: fields.Datetime
					DeletedAt: fields.Datetime

					email:    fields.Email
					username: fields.String
					password: fields.Password
					verified: fields.Bool
					active:   fields.Bool

					persona: fields.Enum & {
						Vals: ["guest", "user", "admin", "owner"]
						Default: "user"
					}
				}
			}
		}
	}

	# Example Usage   (dm is short for datamodel)

	  $ hof dm list   (print known data models)
	  NAME         TYPE       VERSION  STATUS  ID
	  Config       object     -        ok      Config
	  MyDatamodel  datamodel  -        ok      datamodel-abc123

	  $ hof dm tree   (print the structure of the datamodels)

	  $ hof dm diff   (prints a tree based diff of the datamodel)

	  $ hof dm checkpoint -m "a message about this checkpoint"

	  $ hof dm log    (prints the log of changes from latest to oldest)

	  You can also use the -d & -e flags to subselect datamodels and nested values
	
	# Learn more:
	  - https://docs.hofstadter.io/getting-started/data-layer/
	  - https://docs.hofstadter.io/data-modeling/
	"""
