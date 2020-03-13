package gen

import (
  "text/mustache"

  "github.com/hofstadter-io/dsl-cli/schema"
  "github.com/hofstadter-io/dsl-cli/templates"
)

CommandGen : {
  In: {
    CLI: schema.Cli
    CMD: schema.Command
  }
  Template: templates.CommandTemplate
  Filename: "commands/\(In.CMD.Name).go"
  Out: mustache.Render(Template, In)
}

