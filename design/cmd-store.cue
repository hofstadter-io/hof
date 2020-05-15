package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#StoreCommand: schema.#Command & {
  Name:    "store"
  Usage:   "store"
  Aliases: ["s"]
  Short:   "create, checkpoint, and migrate your storage engines"
  Long:    Short

  OmitRun: true

  Commands: [...schema.#Command] & [
    {
      Name:  "run"
      Usage: "run"
      Aliases: ["local"]
      Short: "run local datastore servers"
      Long:  Short

      Args: [...schema.#Arg] & [
        {
          Name:     "dstype"
          Type:     "string"
          Required: true
          Help:     "datastore type"
        },
        {
          Name:     "name"
          Type:     "string"
          Required: true
          Help:     "datastore name"
        },
      ]
    },
    {
      Name:  "conn"
      Usage: "conn"
      Short: "connect to the local datastore"
      Long:  Short
      Args: [...schema.#Arg] & [
        {
          Name:     "name"
          Type:     "string"
          Required: true
          Help:     "datastore name"
        },
      ]
    },

  ]
},
