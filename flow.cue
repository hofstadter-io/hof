package hof

_force: bool | *false @tag(force,type=bool)
_print: bool | *false @tag(print,type=bool)

watchBuild: {
	@flow(watch/build)

	// have to localize this task in a flow for it to work
	RR:   RepoRoot
	root: RR.out
	dirs: ["cmd", "flow", "lib", "gen"]

	watch: {
		@task(os.Watch)
		globs: [ for d in dirs {"\(root)/\(d)/**/*.go"}]
		handler: {
			event?: _
			compile: {
				@task(os.Exec)
				cmd: ["go", "install", "\(root)/cmd/hof"]
				env: {
					CGO_ENABLE: "0"
				}
				exitcode: _
			}
			now: {
				dep: compile.exitcode
				n:   string @task(gen.Now)
				s:   "\(n) (\(dep))"
			}
			alert: {
				@task(os.Stdout)
				dep:  now.s
				text: "hof rebuilt \(now.s)\n"
			}
		}
	}
}

_flow: {
  diff: GitDiff & {
		if _print {
			@print()
		}
		ref: "_dev"
	}
  _shouldi: ShouldI & { files: diff.files }
}

_cond: {
	shouldi: yes: bool | *false
	stdout: string
	stderr: string
	if _force || shouldi.yes {
		@task(os.Exec)
		if _print {
			@print()
		}
	}
}

diff: F=_flow & {
	@flow(diff)
	print: {
		@task(os.Stdout)
		text: F.diff.stdout
	}
}

build: F= _flow & {
	@flow(build)

	gen: {
		_cond
    shouldi: F._shouldi & { globs: ["design/"] }
		run: "hof gen hof.cue"
	}

	cli: {
		_cond
    shouldi: F._shouldi & { globs: ["go.*", "cmd/", "flow", "lib/", "schema/", "script/"] }
		run: "go install ./cmd/hof"
	}

  docs: {
		[string]: {
			dir: "docs"
			#after: { $cli: F.cli }
			_cond
		}

		schemas: {
			shouldi: F._shouldi & { globs: ["schema/"] }
			run: "make schemas"
		}

		gen: {
			shouldi: F._shouldi & { globs: ["docs/", "schema/", "flow/tasks/*/*.cue"] }
			run: "make gen"
		}

		cmdhelp: {
			shouldi: F._shouldi & { globs: ["cmd/"] }
			run: "make cmdhelp"
		}

		highlight: {
			shouldi: F._shouldi & { globs: ["docs/code/"] }
			run: "make highlight"
			#after: { $gen: gen, $schemas: schemas }
		}
  }
}

images: F=_flow & {
	@flow(images)
	_reg: "ghcr.io/hofstadter-io"
	for _,tool in ["black", "csharpier", "prettier"] {
		(tool): {
			_cond
			dir: "formatters/tools/\(tool)"
			shouldi: F._shouldi & { globs: [dir] }
			run: "docker build -t \(_reg)/\(tool):dirty -f Dockerfile.debian ."
		}
	}
}