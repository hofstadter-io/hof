package schema

import (
  "strings"
)

Arg : {
  Name:        string
  argName: strings.ToCamel(Name)
  ArgName: strings.ToTitle(Name)

  Type:        string

  // "this".Type
  // Default:     Type

  Help:        string
  Validation:  string
}

