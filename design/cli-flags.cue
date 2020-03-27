package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

CliPflags :: [...schema.Flag] & [
  {
    Name:    "config"
    Type:    "string"
    Default: ""
    Help:    "Some config file path"
    Long:    "config"
    Short:   "c"
  },
  {
    Name:    "identity"
    Long:    "identity"
    Short:   "I"
    Type:    "string"
    Default: ""
    Help:    "the Studios Auth Identity to use during this hof execution"
  },
  {
    Name:    "context"
    Long:    "context"
    Short:   "C"
    Type:    "string"
    Default: ""
    Help:    "the Studios Context to use during this hof execution"
  },
  {
    Name:    "account"
    Long:    "account"
    Short:   "A"
    Type:    "string"
    Default: ""
    Help:    "the Studios Account to use during this hof execution"
  },
  {
    Name:    "project"
    Long:    "project"
    Short:   "P"
    Type:    "string"
    Default: ""
    Help:    "the Studios Project to use during this hof execution"
  },
]

