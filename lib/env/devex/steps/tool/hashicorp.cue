package tool

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

hashicorp: terraform: _hashicorpBin & {#tool: "terraform"}
hashicorp: packer: _hashicorpBin & {#tool: "packer"}

_hashicorpBin: env.Exec & {
	#ver:  string | *"1.14.3"
	#arch: *"arm64" | "amd64"

	#tool: string

	args: ["sh", "-c", _script]
	_file:   "\(#tool)_\(#ver)_linux_\(#arch).zip"
	_src:    "https://releases.hashicorp.com/\(#tool)/\(#ver)/\(_file)"
	_script: """
  set -eou pipefail

  cd /tmp
  wget -q \(_src)
  unzip \(_file)
  mv \(#tool) /usr/local/bin/\(#tool)
  rm -rf /tmp/*
  """
}
