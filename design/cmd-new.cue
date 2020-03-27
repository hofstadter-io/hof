package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

NewCommand :: schema.Command & {
  Name:  "new"
  Usage: "new"
  Short: "create a new project or subcomponent or files"
  Long:  Short + ", depending on the context"
}
