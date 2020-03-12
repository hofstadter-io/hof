package gen

import (
  "text/mustache"

  "github.com/hofstadter-io/dsl-cli/schema"
  "github.com/hofstadter-io/dsl-cli/templates"
)

TestGen : {
  _In: {
    CLI: schema.Cli
  }
  _Template: templates.TestTemplate
  _Filename: "main.go"
  _Out: mustache.Render(_Template, _In)
}

