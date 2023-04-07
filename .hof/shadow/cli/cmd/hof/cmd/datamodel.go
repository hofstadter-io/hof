package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/datamodel"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var datamodelLong = `Data models are values or objects used in many hof processes and modules.
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

  $ hof dm info   (print the structure of the datamodels)

  $ hof dm diff   (prints a tree based diff of the datamodel)

  $ hof dm checkpoint -m "a message about this checkpoint"

  $ hof dm log    (prints the log of changes from latest to oldest)

  You can also use the -d & -e flags to subselect datamodels and nested values

# Learn more:
  - https://docs.hofstadter.io/getting-started/data-layer/
  - https://docs.hofstadter.io/data-modeling/`

func init() {

	DatamodelCmd.PersistentFlags().StringSliceVarP(&(flags.DatamodelPflags.Datamodels), "datamodel", "d", nil, "specify one or more datamodels to operate on")
	DatamodelCmd.PersistentFlags().StringSliceVarP(&(flags.DatamodelPflags.Expression), "expr", "e", nil, "CUE paths to select outputs, depending on the command")
}

var DatamodelCmd = &cobra.Command{

	Use: "datamodel",

	Aliases: []string{
		"dm",
	},

	Short: "manage, diff, and migrate your data models",

	Long: datamodelLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := DatamodelCmd.HelpFunc()
	ousage := DatamodelCmd.UsageFunc()
	help := func(cmd *cobra.Command, args []string) {
		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	DatamodelCmd.SetHelpFunc(thelp)
	DatamodelCmd.SetUsageFunc(tusage)

	DatamodelCmd.AddCommand(cmddatamodel.CheckpointCmd)
	DatamodelCmd.AddCommand(cmddatamodel.DiffCmd)
	DatamodelCmd.AddCommand(cmddatamodel.InfoCmd)
	DatamodelCmd.AddCommand(cmddatamodel.ListCmd)
	DatamodelCmd.AddCommand(cmddatamodel.LogCmd)

}
