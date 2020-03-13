package templates

CommandTemplate : FakeCommandTemplate

FakeCommandTemplate : "CommandTemplate"

RealCommandTemplate : """
{{> package-name.go CTX=RepeatedContext}}

{{#with RepeatedContext as |RC| }}
{{#with DslContext as |CLI| }}

import (
	{{#if (or RC.omit-run RC.body)}}
	{{else}}
	"fmt"
	{{/if}}

	// custom imports
	{{#each RC.imports as |I|}}
	{{I.as}} "{{{ I.path }}}"
	{{/each}}

	// infered imports
	{{#dotpath "args.required" RC false}}
	{{#with . as |D|}}
	{{#if (contains D "Error")}}
	{{else}}
	"os"
	{{/if}}
	{{/with}}
	{{/dotpath}}

	{{#if RC.flags}}
	"github.com/spf13/viper"
	{{else}}
		{{#if RC.pflags}}
	"github.com/spf13/viper"
		{{/if}}
	{{/if}}
	"github.com/spf13/cobra"

	{{#if RC.commands}}
	"{{CLI.package}}/{{lower RC.name}}"
	{{/if}}
)

// Tool:   {{CLI.name}}
// Name:   {{RC.name}}
// Usage:  {{{RC.usage}}}
// Parent: {{{RC.parent}}}


{{#if RC.long}}
var {{camelT RC.name}}Long = `{{{RC.long}}}`
{{/if}}

{{> "flag-var.go" RC }}

{{> "flag-init.go" RC }}

var {{camelT RC.name}}Cmd = &cobra.Command {
	{{#if RC.hidden}}
	Hidden: true,
	{{/if}}

	{{#if RC.usage}}
	Use: "{{{RC.usage}}}",
	{{else}}
	Use: "{{{RC.name}}}",
	{{/if}}

	{{#if RC.aliases}}
	Aliases: []string{
		{{#each RC.aliases as |AL|}}"{{AL}}",
		{{/each}}
	},
	{{/if}}

	{{#if RC.short}}
	Short: "{{{RC.short}}}",
	{{/if}}

	{{#if RC.long}}
	Long: {{camelT RC.name}}Long,
	{{/if}}

	{{#if RC.persistent-prerun}}
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug("In PersistentPreRun {{RC.name}}Cmd", "args", args)
		{{> args-parse.go RC }}

		{{#if RC.persistent-prerun-body}}
		{{{ RC.persistent-prerun-body}}}
		{{/if}}
	},
	{{/if}}
	{{#if RC.prerun}}
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Debug("In PreRun {{RC.name}}Cmd", "args", args)
		{{> args-parse.go RC }}

		{{#if RC.prerun-body}}
		{{{ RC.prerun-body}}}
		{{/if}}
	},
	{{/if}}
	{{#unless RC.omit-run}}
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In {{RC.name}}Cmd", "args", args)
		{{> args-parse.go RC }}

		{{#if RC.body}}
		{{{ RC.body}}}
		{{else}}
		fmt.Println("{{replace RC.pkgPath "/" " " -1}}:", {{#each RC.args}}
		{{camel name}},
		{{/each}})
		{{/if}}
	},
	{{/unless}}
	{{#if RC.persistent-postrun}}
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logger.Debug("In PersistentPostRun {{RC.name}}Cmd", "args", args)
		{{> args-parse.go RC }}

		{{#if RC.persistent-postrun-body}}
		{{{ RC.persistent-postrun-body}}}
		{{/if}}
	},
	{{/if}}
	{{#if RC.postrun}}
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Debug("In PostRun {{RC.name}}Cmd", "args", args)
		{{> args-parse.go RC }}

		{{#if RC.postrun-body}}
		{{{ RC.postrun-body}}}
		{{/if}}
	},
	{{/if}}
}


{{#if (eq RC.parent CLI.name) }}
func init() {
	RootCmd.AddCommand({{camelT RC.name}}Cmd)
}
{{/if}}

{{#if commands}}
func init() {
	// add sub-commands to this command when present

	{{#each RC.commands as |C|}}
	{{camelT RC.name}}Cmd.AddCommand({{lower RC.name}}.{{camelT C.name}}Cmd)
	{{/each}}
}
{{/if}}

{{/with}}
{{/with}}

"""
