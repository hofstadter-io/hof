package templates

import (
  "github.com/hofstadter-io/dsl-cli/partials"
)

CommandTemplate : partials.AllPartials + RealCommandTemplate

RealCommandTemplate : """
{{ if .CMD.Parent }}
package {{ .CMD.Parent.Name }}
{{ else }}
package commands
{{ end }}

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  {{ if or .CMD.Flags .CMD.Pflags }}
  "github.com/spf13/viper"
  {{ end }}

  {{ if .CMD.Imports }}
	{{ range $i, $I := .CMD.imports }}
	{{ $I.As }} "{{ $I.Path }}"
	{{ end }}
	{{ end }}

	{{ if .CMD.Commands }}
	"{{ .CLI.Package }}/{{ .CMD.cmdName }}"
	{{ end }}
)

{{ if .CMD.Long }}
var {{ .CMD.Name }}Long = `{{ .CMD.Long }}`
{{ end }}

{{ template "flag-vars" .CMD }}
{{ template "flag-init" .CMD }}

var {{ .CMD.CmdName }}Cmd = &cobra.Command{

  {{ if .CMD.Usage}}
  Use: "{{ .CMD.Usage }}",
  {{ else }}
  Use: "{{ .CMD.Name }}",
  {{ end }}

	{{ if .CMD.Hidden }}
	Hidden: true,
	{{ end }}

	{{ if .CMD.Aliases }}
	Aliases: []string{
		{{range $i, $AL := .CMD.Aliases}}"{{$AL}}",
		{{end}}
	},
	{{ end }}

  {{ if .CMD.Short}}
  Short: "{{ .CMD.Short }}",
  {{ end }}

  {{ if .CMD.Long }}
  Long: {{ .CMD.Name }}Long,
  {{ end }}

  {{ if .CMD.PersistentPrerun }}
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    {{ template "args-parse" .CMD }}

    {{ if .CMD.PersistentPrerunBody }}
    {{ .CMD.PersistentPrerunBody }}
    {{ end }}
  },
  {{ end }}

  {{ if .CMD.Prerun }}
  PreRun: func(cmd *cobra.Command, args []string) {
    {{ template "args-parse" .CMD }}

    {{ if .CMD.PrerunBody }}
    {{ .CMD.PrerunBody }}
    {{ end }}
  },
  {{ end }}

  {{ if not .CMD.OmitRun}}
  Run: func(cmd *cobra.Command, args []string) {
    {{ template "args-parse" .CMD }}

    {{ if .CMD.Body}}
    {{ .CMD.Body}}
    {{ end }}
  },
  {{ end }}

  {{ if .CMD.PersistentPostrun}}
  PersistentPostRun: func(cmd *cobra.Command, args []string) {
    {{ template "args-parse" .CMD }}

    {{ if .CMD.PersistentPostrunBody}}
    {{ .CMD.PersistentPostrunBody}}
    {{ end }}
  },
  {{ end }}

  {{ if .CMD.Postrun}}
  PostRun: func(cmd *cobra.Command, args []string) {
    {{ template "args-parse" .CMD }}

    {{ if .CMD.PostrunBody }}
    {{ .CMD.PostrunBody }}
    {{ end }}
  },
  {{ end }}
}

{{if .CMD.Commands}}
func init() {
	{{ range $i, $C := .CMD.Commands -}}
  {{ $.CMD.CmdName }}Cmd.AddCommand({{ $.CMD.cmdName }}.{{ $C.CmdName }}Cmd)
	{{ end -}}
}
{{ end }}

"""
