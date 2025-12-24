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
    set -eux
    apt-get update -y
    apt-get install -y --no-install-recommends \(strings.Join(#pkgs, " "))
    apt-get dist-clean
    rm -rf /var/lib/apt/lists/*
    """
    args: ["bash", "-c", _script]
  }

  // 
}
