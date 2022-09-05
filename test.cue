package hof

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

	hack: {
		test: string | *"TestMainFlow" @tag(test)
		@flow(test/hack)
		prt: {text: "testing: \(test)\n"} @task(os.Stdout)
		run: {
			@task(os.Exec)
			cmd: ["bash", "-c", script]
			dir:    "flow"
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

	gen: {
		@flow(test/gen)
		run: BashTest & {
			dir: "test/templates"
			script: """
				set -e
				# fetch CUE deps
				hof mod vendor cue
				# generate templates
				hof gen
				# should have no diff
				git diff --exit-code
				"""
		}
	}

	render: {
		@flow(test/render)
		run: GoTest & {
			dir: "test/render"
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
}
