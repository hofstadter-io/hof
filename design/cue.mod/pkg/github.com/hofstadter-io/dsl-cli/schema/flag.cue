package schema

import (
  "strings"
)

Flag : {
  Name:    string
  flagName: strings.ToCamel(Name)
  FlagName: strings.ToTitle(Name)

  Type:        string
  Default:     string
  Help:        string
  Long:        string
  Short:       string
  Validation:  string
}

