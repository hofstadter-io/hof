package utils

import "strings"

RepoRoot: { 
  $id: "tool/exec.Run"
  cmd: ["bash", "-c", "git rev-parse --show-toplevel"]
  stdout: string
  Out: strings.TrimSpace(stdout)
}
