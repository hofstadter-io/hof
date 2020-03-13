package gen

import (
  "text/mustache"

  "github.com/hofstadter-io/dsl-cli/schema"
  "github.com/hofstadter-io/dsl-cli/templates"
)

RootGen : {
  In: {
    CLI: schema.Cli
  }
  Template: templates.RootTemplate
  Filename: "commands/root.go"
  Out: mustache.Render(Template, In)
}

