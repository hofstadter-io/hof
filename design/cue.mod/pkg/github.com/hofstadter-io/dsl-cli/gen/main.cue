package gen

import (
  "text/mustache"

  "github.com/hofstadter-io/dsl-cli/schema"
  "github.com/hofstadter-io/dsl-cli/templates"
)

MainGen : {
  In: {
    CLI: schema.Cli
  }
  Template: templates.MainTemplate
  Filename: "main.go"
  Out: mustache.Render(Template, In)
}

