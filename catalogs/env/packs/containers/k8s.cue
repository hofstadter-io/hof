@experiment(aliasv2)
package containers

import (
	"strings"

	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/schemas/env"
)

// basically the same as github, but without the version in teh filname
crane: {
	_name: "go-containerregistry"
	_repo: "google/\(_name)"

	// we could consider normalizing across all the github projects we download from
	#ver:    string | *"0.20.7"
	#arch:   string | *"arm64" | "x86_64"
	#distro: *"Linux" | "Darwin" | "Windows"

	#bins: ["crane", "gcrane", "krane"]
	_bins: strings.Join(#bins, " ")
	_file: "\(_name)_\(#distro)_\(#arch).tar.gz"
	_src:  "https://github.com/\(_repo)/releases/download/v\(#ver)/\(_file)"

	files: env.#Dir & {
		@env(crane-files)
		sources: [ctr]
		include: ["work/*"]
		trimPrefix: "work"
	}
	ctr: env.#Container & {
		@env(crane-ctr)
		from: bases.debian13.minimal
		steps: [
			env.Sh & {
				script: """
          wget -q \(_src)
          tar -xzf \(_file)
          rm -rf \(_file)
          """
			},
		]
	}
}
