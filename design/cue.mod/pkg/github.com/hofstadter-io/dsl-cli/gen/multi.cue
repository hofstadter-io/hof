package gen

import (
  "text/mustache"

  "github.com/hofstadter-io/dsl-cli/schema"
  "github.com/hofstadter-io/dsl-cli/templates"
)

MultiGen : {
  _In: {
    CLI: schema.Cli
    CMD: schema.Command
  }
  _Template: templates.MultiTemplate
  _Filename: "commands/\(_In.CMD.Name).go"
  _Out: mustache.Render(_Template, _In)
}

