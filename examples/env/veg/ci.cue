@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

ci: {

	// These are set to developer local by default
	// in CI, these are different git repo checkouts
	// - for a PR, the target and current
	// - for a branch w/o PR, main branch, until bfg
	#prev: _ | *src.repo
	#next: _ | *src.disk

	// this is what we run to verify the cue mod
	cuemod: {
		// container with cuemod files
		ctr: env.#Container & {
			@id(ci-cuemod-with-src)
			from: root.ctr.dev
			steps: [
				env.Dir & {source: root.src.cuemod},
			]
		}
		_runner: env.#Container & {
			#script: string
			from:    ctr
			steps: [
				env.Sh & {script: #script},
			]
		}

		shouldi: env.#Shouldi & {
			@env(ci-cuemod-shouldi)
			changes: env.#Changes & {
				prev: #prev
				next: #next
			}
			include: root.src.cuemod.include
			then: env.#Task & {
				@env(ci-cuemod-task)
				steps: [
					// fmt/lint
					[_runner & {#cmd: "cue fmt ./..."}],

					// vet
					[_runner & {#cmd: "cue vet -c=false ./..."}],
				]
			}
		}

		cmd: env.#Command & {
			name: "ci.cuemod"
			tasks: {
				name: "check"
				steps: [[shouldi]]
			}
		}

	}

}
