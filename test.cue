package hof

import "strings"

BashTest: {
	@task(os.Exec)
	script: string
	cmd: ["bash", "-c", script]
}

GoTest: {
	@task(os.Exec)
	cmd: ["bash", "-c", scripts.DEV]
}

scripts: {
	DEV: string | *"""
		rm -rf .workdir
		go test -cover ./
		"""

	CI: string | *"""
		rm -rf .workdir
		go test -cover ./ -coverprofile cover.out -json > tests.json
		"""
}

watchTest: {
	@flow(watch/test)

	watch: {
		@task(fs.Watch)
		first: true
		globs: [
			"lib/structural/**/*.*",
		]
		handler: {
			event?: _
			compile: {
				@task(os.Exec)
				cmd: ["hof", "flow", "-f", "test/hack"]
			}
		}
	}
}

tests: {
	// want to discover nested too
	// @flow(test)

	// primarily used by devs to select isolated tests to run
	one: {
		@flow(test/one)

		test: string | *"TestModAuthdApikeysTests" @tag(test)
		dir:  string | *"lib/mod"                  @tag(dir)

		prt: {text: "testing: \(test)\n"} @task(os.Stdout)

		run: {
			@task(os.Exec)
			cmd: ["bash", "-c", script]
			"dir":  dir
			script: """
      rm -rf .workdir
      go test -run \(test) . 
      """
		}
	}

	flow: {
		@flow(test/flow)
		run: GoTest & {
			// dir: "lib/flow" // panics, segfault
			dir: "flow"
		}
	}

	dm: {
		@flow(test/dm)
		run: GoTest & {
			dir: "lib/datamodel/test"
		}
	}

	gen: {
		@flow(test/gen)
		run: BashTest & {
			dir: "test/templates"
			script: """
				set -e
				# deps & gen
				hof mod vendor
				hof gen

				# should have no diff
				git diff
				git status
				git diff --exit-code
				"""
		}
	}

	render: {
		@flow(test/render)
		run: GoTest & {
			dir: "test/render"
			script: """
				echo "no-op"
				echo "now run in the dagger workflow"
				"""
		}
	}

	create: {
		@flow(test/create)
		test: BashTest & {
			dir: "test/create"
			script: """
				cd test_01 && make test && cd ..
				cd test_02 && make test && cd ..
				"""
		}
	}

	st: {
		@flow(test/st)
		run: GoTest & {
			dir: "lib/structural"
		}
	}

	mod: {
		@flow(test/mod)
		run: GoTest & {
			dir: "lib/mod"
		}
	}
	fmt: {
		@flow(test/fmt)
		run: GoTest & {
			dir: "formatters/test"
		}
	}
	cue: {
		all: {
			@flow(test/cue)
			run: GoTest & {
				dir: "lib/cuecmd"
			}
		}

		hack: {
			@flow(test/cue/hack)
			run: {
				@task(os.Exec)
				cmd: ["bash", "-c", _script]
				dir: "lib/cuecmd"
				_script: """
					rm -rf .workdir
					go test -cover -run Eval ./
					"""
			}
		}

		for k in ["def", "eval", "export", "vet"] {
			(k): {
				let K = strings.ToTitle(k)
				#hof: Flow: {Root: true, Name: "test/cue/\(k)"}
				run: {
					@task(os.Exec)
					cmd: ["bash", "-c", _script]
					dir:     "lib/cuecmd"
					_script: """
						rm -rf .workdir
						go test -cover -run \(K) ./
						"""
				}
			}
		}
	}
}
