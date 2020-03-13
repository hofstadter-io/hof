package templates

import (
  "github.com/hofstadter-io/dsl-cli/partials"
)

RootTemplate : partials.AllPartials + RealRootTemplate

FakeRootTemplate : "FakeRootTemplate"

RealRootTemplate : """
package commands

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  {{ if or .CLI.Flags .CLI.Pflags }}
  "github.com/spf13/viper"
  {{ end }}

  {{ if .CLI.Imports }}
	{{ range $i, $I := .CLI.imports }}
	{{ $I.As }} "{{ $I.Path }}"
	{{ end }}
	{{ end }}
)

{{ if .CLI.Long }}
var {{ .CLI.Name }}Long = `{{ .CLI.Long }}`
{{ end }}

{{ template "flag-vars" }}
{{ template "flag-init" }}

var (
	RootCmd = &cobra.Command{

		{{ if .CLI.Usage}}
		Use: "{{ .CLI.Usage }}",
		{{ else }}
		Use: "{{ .CLI.Name }}",
		{{ end }}

		{{ if .CLI.Short}}
		Short: "{{ .CLI.Short }}",
		{{ end }}

		{{ if .CLI.Long }}
		Long: {{ .CLI.Name }}Long,
		{{ end }}

		{{ if .CLI.PersistentPrerun }}
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			{{ template "args-parse" .CLI }}

			{{ if .CLI.PersistentPrerunBody }}
			{{ .CLI.PersistentPrerunBody }}
			{{ end }}
		},
		{{ end }}

		{{ if .CLI.Prerun }}
		PreRun: func(cmd *cobra.Command, args []string) {
			{{ template "args-parse" .CLI }}

			{{ if .CLI.PrerunBody }}
			{{ .CLI.PrerunBody }}
			{{ end }}
		},
		{{ end }}

		{{ if not .CLI.OmitRun}}
		Run: func(cmd *cobra.Command, args []string) {
			{{ template "args-parse" .CLI }}

			{{ if .CLI.Body}}
			{{ .CLI.Body}}
			{{ end }}
		},
		{{ end }}

		{{ if .CLI.PersistentPostrun}}
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			{{ template "args-parse" .CLI }}

			{{ if .CLI.PersistentPostrunBody}}
			{{ .CLI.PersistentPostrunBody}}
			{{ end }}
		},
		{{ end }}

		{{ if .CLI.Postrun}}
		PostRun: func(cmd *cobra.Command, args []string) {
			{{ template "args-parse" .CLI }}

			{{ if .CLI.PostrunBody }}
			{{ .CLI.PostrunBody }}
			{{ end }}
		},
		{{ end }}
	}
)

{{if .CLI.Commands}}
func init() {
	{{ range $i, $C := .CLI.Commands }}
	RootCmd.AddCommand({{ $C.Name }}Cmd)
	{{ end }}
}
{{ end }}

"""

RootRootTemplate : """
{{#with DslContext as |CLI| }}
package commands

import (
	{{#if (or CLI.omit-run CLI.body)}}
	{{else}}
	"fmt"
	{{/if}}

	// custom imports
	{{#each CLI.imports as |I|}}
	{{I.as}} "{{{ I.path }}}"
	{{/each}}

	// infered imports
	{{#dotpath "args.required" CLI false}}
	{{#with . as |D|}}
	{{#if (contains D "Error")}}
	{{else}}
	"os"
	{{/if}}
	{{/with}}
	{{/dotpath}}

	{{#if CLI.flags}}
	"github.com/spf13/viper"
	{{else}}
		{{#if CLI.pflags}}
	"github.com/spf13/viper"
		{{/if}}
	{{/if}}
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
)

{{#if CLI.long}}
var {{camelT CLI.name}}Long = `{{{CLI.long}}}`
{{/if}}

{{> "flag-var.go" CLI }}

{{> "flag-init.go" CLI }}

var (
	RootCmd = &cobra.Command{

		{{#if CLI.usage}}
		Use: "{{{CLI.usage}}}",
		{{else}}
		Use: "{{{CLI.name}}}",
		{{/if}}

		{{#if CLI.short}}
		Short: "{{{CLI.short}}}",
		{{/if}}

		{{#if CLI.long}}
		Long: {{camelT CLI.name}}Long,
		{{/if}}

		{{#if CLI.persistent-prerun}}
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			{{> args-parse.go CLI }}

			{{#if CLI.persistent-prerun-body}}
			{{{ CLI.persistent-prerun-body}}}
			{{/if}}
		},
		{{/if}}
		{{#if CLI.prerun}}
		PreRun: func(cmd *cobra.Command, args []string) {
			logger.Debug("In PreRun {{CLI.name}}Cmd", "args", args)
			{{> args-parse.go CLI }}

			{{#if CLI.prerun-body}}
			{{{ CLI.prerun-body}}}
			{{/if}}
		},
		{{/if}}
		{{#unless CLI.omit-run}}
		Run: func(cmd *cobra.Command, args []string) {
			logger.Debug("In {{CLI.name}}Cmd", "args", args)
			{{> args-parse.go CLI }}

			{{#if CLI.body}}
			{{{ CLI.body}}}
			{{else}}
			fmt.Println("{{replace CLI.pkgPath "/" " " -1}}:", {{#each CLI.args}}
			{{camel name}},
			{{/each}})
			{{/if}}
		},
		{{/unless}}
		{{#if CLI.persistent-postrun}}
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			logger.Debug("In PersistentPostRun {{CLI.name}}Cmd", "args", args)
			{{> args-parse.go CLI }}

			{{#if CLI.persistent-postrun-body}}
			{{{ CLI.persistent-postrun-body}}}
			{{/if}}
		},
		{{/if}}
		{{#if CLI.postrun}}
		PostRun: func(cmd *cobra.Command, args []string) {
			logger.Debug("In PostRun {{CLI.name}}Cmd", "args", args)
			{{> args-parse.go CLI }}

			{{#if CLI.postrun-body}}
			{{{ CLI.postrun-body}}}
			{{/if}}
		},
		{{/if}}
	}
)

{{/with}}
"""

