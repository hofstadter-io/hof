package gen

import (
  "text/template"

  "github.com/hofstadter-io/dsl-cli/schema"
  "github.com/hofstadter-io/dsl-cli/templates"
)

RootGen : {
  In: {
    CLI: schema.Cli
  }
  Template: templates.RootTemplate
  Filename: "commands/root.go"
  Out: template.Execute(Template, In)
}

