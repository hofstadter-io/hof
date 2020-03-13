package templates

import (
  "github.com/hofstadter-io/dsl-cli/partials"
)

RootTemplate : partials.AllPartials + RealRootTemplate

RealRootTemplate : """
package commands
// {{ .CLI.Package }}

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

{{ template "flag-vars" .CLI }}
{{ template "flag-init" .CLI }}

var RootCmd = &cobra.Command{

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

{{if .CLI.Commands}}
func init() {
	{{ range $i, $C := .CLI.Commands }}
	RootCmd.AddCommand({{ $C.CmdName }}Cmd)
	{{ end }}
}
{{ end }}

"""
