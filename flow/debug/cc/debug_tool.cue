package debug

import "tool/cli"
import "tool/os"

command: {
  do: {
    get: os.Getenv & {
      FOO: string
    }
    say: cli.Print & {
      text: "(cc): the cow goes: \(get.FOO)"
    }
  }
}
