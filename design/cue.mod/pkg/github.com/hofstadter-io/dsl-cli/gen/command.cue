package gen

import (
  "text/template"

  "github.com/hofstadter-io/dsl-cli/schema"
  "github.com/hofstadter-io/dsl-cli/templates"
)

CommandGen : {
  In: {
    CLI: schema.Cli
    CMD: schema.Command
  }
  Template: templates.CommandTemplate
  if In.CMD.Parent == _|_ {
    Filename: "commands/\(In.CMD.Name).go"
  } 
  if In.CMD.Parent != _|_ {
    Filename: "commands/\(In.CMD.Parent.Name)/\(In.CMD.Name).go"
  }
  Out: template.Execute(Template, In)
}

