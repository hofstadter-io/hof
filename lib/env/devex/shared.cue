package devex

import (
  "strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

_steps: {

  // generalized apt package install
  apt: env.Exec & {
    #pkgs: [...string]
    _script: """
    apt-get update -y
    apt-get install -y \(strings.Join(#pkgs, " "))
    rm -rf /var/lib/apt/lists/*
    """
    args: ["bash", "-c", _script]
  }

  // 
}
