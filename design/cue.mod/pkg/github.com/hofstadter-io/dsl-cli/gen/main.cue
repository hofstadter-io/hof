package gen

import (
  "text/template"

  "github.com/hofstadter-io/dsl-cli/schema"
  "github.com/hofstadter-io/dsl-cli/templates"
)

MainGen : {
  In: {
    CLI: schema.Cli
  }
  Template: templates.MainTemplate
  Filename: "main.go"
  Out: template.Execute(Template, In)
}

