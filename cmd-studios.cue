package hof

import (
	"github.com/hofstadter-io/cuemod--cli-golang/schema"
)

StudiosCommand :: schema.Command & {
  Name:  "studios"
  Usage: "studios"
  Aliases: ["s"]
  Short: "commands for working with Hofstadter Studios"
  Long: """
    Hofstadter Studios makes it easy to develop and launch both
    hof-lang modules as well as pretty much any code or application
  """

  OmitRun: true

  Commands: [
    studiosAppCmd,
    studiosDatabaseCmd,
    studiosContainerCmd,
    studiosFunctionCmd,
    studiosConfigCmd,
    studiosSecretCmd,
  ]

}

studiosAppCmd :: schema.Command & {

  Name:  "app"
  Usage: "app"
  Short: "Work with Hofstadter Studios apps"
  Long:  "Work with Hofstadter Studios apps"

  Commands: [
    schema.Command & {
      Name:  "list"
      Usage: "list"
      Short: "List your apps"
      Long:  "List your Studios apps"
    },
    schema.Command & {
      Name:  "get"
      Usage: "get <name or id>"
      Short: "Get a Studios app"
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "create"
      Usage: "create <name> [input]"
      Short: "Create a Studios app"
      Long:  "Create a Studios app by name with extra creation values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "update"
      Usage: "update <name> <input>"
      Short: "Update a Studios app"
      Long:  "Update a Studios app by name with extra update values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "deploy"
      Usage: "deploy <name> <input>"
      Short: "Deploy a Studios app"
      Long:  "Deploy a Studios app by name with extra update values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "status"
      Usage: "status <name or id>"
      Short: "Get the status of a Studios app."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "push"
      Usage: "push <name or id>"
      Short: "Push a Studios app."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "pull"
      Usage: "pull <name or id>"
      Short: "Pull a Studios app."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "reset"
      Usage: "reset <name or id>"
      Short: "Reset a Studios app."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "shutdown"
      Usage: "shutdown <name or id>"
      Short: "Shutdown a Studios app."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "delete"
      Usage: "delete <name or id>"
      Short: "Delete a Studios app."
      Long:  Short
      Args: [ identArg]
    },
  ]

  // Leave open for parents annotation
  ...
}

studiosDatabaseCmd :: schema.Command & {

  Name:  "database"
  Usage: "database"
  Short: "Work with Hofstadter Studios databases"
  Long:  "Work with Hofstadter Studios databases"

  Commands: [
    schema.Command & {
      Name:  "list"
      Usage: "list"
      Short: "List your databases"
      Long:  "List your Studios databases"
    },
    schema.Command & {
      Name:  "get"
      Usage: "get <name or id>"
      Short: "Get a Studios database"
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "create"
      Usage: "create <name> [input]"
      Short: "Create a Studios database"
      Long:  "Create a Studios database by name with extra creation values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "update"
      Usage: "update <name> <input>"
      Short: "Update a Studios database"
      Long:  "Update a Studios database by name with extra update values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "status"
      Usage: "status <name or id>"
      Short: "Get the status of a Studios database."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "save"
      Usage: "save <name or id> <backup-name>"
      Short: "Save a Studios database under a named reference."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "restore"
      Usage: "restore <name or id> <backup-name>"
      Short: "Restore a Studios database."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "delete"
      Usage: "delete <name or id>"
      Short: "Delete a Studios database."
      Long:  Short
      Args: [ identArg]
    },
  ]
}

studiosContainerCmd :: schema.Command & {

  Name:  "container"
  Usage: "container"
  Aliases: ["cont"]
  Short: "Work with Hofstadter Studios containers"
  Long:  "Work with Hofstadter Studios containers"

  Commands: [
    schema.Command & {
      Name:  "call"
      Usage: "call"
      Short: "Call a container"
      Long:  "Call your Studios container"
    },

    schema.Command & {
      Name:  "list"
      Usage: "list"
      Short: "List your containers"
      Long:  "List your Studios containers"
    },
    schema.Command & {
      Name:  "get"
      Usage: "get <name or id>"
      Short: "Get a Studios container"
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "create"
      Usage: "create <name> [input]"
      Short: "Create a Studios container"
      Long:  "Create a Studios container by name with extra creation values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "update"
      Usage: "update <name> <input>"
      Short: "Update a Studios container"
      Long:  "Update a Studios container by name with extra update values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "deploy"
      Usage: "deploy <name> <input>"
      Short: "Deploy a Studios container"
      Long:  "Deploy a Studios container by name with extra update values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "status"
      Usage: "status <name or id>"
      Short: "Get the status of a Studios container."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "push"
      Usage: "push <name or id>"
      Short: "Push a Studios container."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "pull"
      Usage: "pull <name or id>"
      Short: "Pull a Studios container."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "reset"
      Usage: "reset <name or id>"
      Short: "Reset a Studios container."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "shutdown"
      Usage: "shutdown <name or id>"
      Short: "Shutdown a Studios container."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "delete"
      Usage: "delete <name or id>"
      Short: "Delete a Studios container."
      Long:  Short
      Args: [ identArg]
    },
  ]

  // Leave open for parents annotation
  ...
}


studiosFunctionCmd :: schema.Command & {

  Name:  "function"
  Usage: "function"
  Aliases: ["cont"]
  Short: "Work with Hofstadter Studios functions"
  Long:  "Work with Hofstadter Studios functions"

  Commands: [
    schema.Command & {
      Name:  "call"
      Usage: "call"
      Short: "Call a function"
      Long:  "Call your Studios function"
    },

    schema.Command & {
      Name:  "list"
      Usage: "list"
      Short: "List your functions"
      Long:  "List your Studios functions"
    },
    schema.Command & {
      Name:  "get"
      Usage: "get <name or id>"
      Short: "Get a Studios function"
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "create"
      Usage: "create <name> [input]"
      Short: "Create a Studios function"
      Long:  "Create a Studios function by name with extra creation values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "update"
      Usage: "update <name> <input>"
      Short: "Update a Studios function"
      Long:  "Update a Studios function by name with extra update values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "deploy"
      Usage: "deploy <name> <input>"
      Short: "Deploy a Studios function"
      Long:  "Deploy a Studios function by name with extra update values as input"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "status"
      Usage: "status <name or id>"
      Short: "Get the status of a Studios function."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "push"
      Usage: "push <name or id>"
      Short: "Push a Studios function."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "pull"
      Usage: "pull <name or id>"
      Short: "Pull a Studios function."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "reset"
      Usage: "reset <name or id>"
      Short: "Reset a Studios function."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "shutdown"
      Usage: "shutdown <name or id>"
      Short: "Shutdown a Studios function."
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "delete"
      Usage: "delete <name or id>"
      Short: "Delete a Studios function."
      Long:  Short
      Args: [ identArg]
    },
  ]

  // Leave open for parents annotation
  ...
}


studiosConfigCmd :: schema.Command & {

  Name:  "config"
  Usage: "config"
  Aliases: [
    "cfg",
  ]
  Short: "Work with Hofstadter Studios configs"
  Long:  "Work with Hofstadter Studios configs"

  Commands: [
    schema.Command & {
      Name:  "list"
      Usage: "list"
      Short: "List your configs"
      Long:  "List your Studios configs"
    },
    schema.Command & {
      Name:  "get"
      Usage: "get <name or id>"
      Short: "Get a Studios config"
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "create"
      Usage: "create <name> <input>"
      Short: "Create a Studios config"
      Long:  "Create a Studios config from a file or key/val pairs"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "update"
      Usage: "update <name> <input>"
      Short: "Update a Studios config"
      Long:  "Update a Studios config from a file or key/val pairs"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "delete"
      Usage: "delete <name or id>"
      Short: "Delete a Studios config. Must not be in use"
      Long:  Short
      Args: [ identArg]
    },
  ]

  // Leave open for parents annotation
  ...
}

studiosSecretCmd :: schema.Command & {

  Name:  "secret"
  Usage: "secret"
  Aliases: [
    "secrets",
    "shh",
  ]
  Short: "Work with Hofstadter Studios secrets"
  Long:  "Work with Hofstadter Studios secrets"

  Commands: [
    schema.Command & {
      Name:  "list"
      Usage: "list"
      Short: "List your secrets"
      Long:  "List your Studios secrets"
    },
    schema.Command & {
      Name:  "get"
      Usage: "get <name or id>"
      Short: "Get a Studios secret"
      Long:  Short
      Args: [ identArg]
    },
    schema.Command & {
      Name:  "create"
      Usage: "create <name> <input>"
      Short: "Create a Studios secret"
      Long:  "Create a Studios secret from a file or key/val pairs"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "update"
      Usage: "update <name> <input>"
      Short: "Update a Studios secret"
      Long:  "Update a Studios secret from a file or key/val pairs"
      Args: [
        nameArg,
        inputArg & {
          Help: "@file or key=val,key2=val2,..."
        },
      ]
    },
    schema.Command & {
      Name:  "delete"
      Usage: "delete <name or id>"
      Short: "Delete a Studios secret. Must not be in use"
      Long:  Short
      Args: [ identArg]
    },
  ]

  // Leave open for parents annotation
  ...
}
