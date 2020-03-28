package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

RunCommand :: schema.Command & {
  Name:  "run"
  Usage: "run [flags] [cmd] [args]"
  Short: "run commands defined by HofCmd"
  Long: "run commands defined by HofCmd. Falls back on cue commands, which also falls back to their own run system"

  Imports: [
    {Path: "fmt"},
    {Path: "os"},
    {Path: CLI.Package + "/lib"},
  ]

  Body: """
    flags := []string{}
    msg, err := lib.Cmd(flags, args, "")
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    fmt.Println(msg)
  """
}
