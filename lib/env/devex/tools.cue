package devex

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

_tools: {
  cue: env.File & {
    path: "/cue"
    source: env.Container & {
      
    }
  }

  cueDL: env.Exec & {
    #ver: string | *"0.15.1"
    #arch: *"arm64" | "amd64"
    args: ["sh", "-c", _script]
    _script: """
      wget https://github.com/cue-lang/cue/releases/download/v0.15.1/cue_v0.15.1_linux_arm64.tar.gz
      tar -xf 
      """
  }

  zsh: env.Exec & {
    args: ["sh", "-c", _script]
    _script: """
      sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
      sed -i 's/^ZSH_THEME=.*/ZSH_THEME="ys"/' /root/.zshrc
      """
  }
}
