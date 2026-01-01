package tool

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

hashicorp: terraform: _hashicorpBin & {#tool: "terraform"}
hashicorp: packer: _hashicorpBin & {#tool: "packer"}

// equivalent to WithFile
_hashicorpBin: env.File & {
	// params
	#ver:  string | *"1.14.3"
	#arch: *"arm64" | "amd64"
	#tool: string

	// internal
	_file:   "\(#tool)_\(#ver)_linux_\(#arch).zip"
	_src:    "https://releases.hashicorp.com/\(#tool)/\(#ver)/\(_file)"

	// spec
	path: #tool
	content: env.#File & {
		path: #tool
		source: env.#Container & {
			from: 
			steps: [env.Bash & {script: "wget -q \(_src) && unzip \(_file)"}]
		}
	}
}
