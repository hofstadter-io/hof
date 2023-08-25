package hof

import "strings"

RepoRoot: {
	@task(os.Exec)
	cmd: ["bash", "-c", "git rev-parse --show-toplevel"]
	stdout: string
	out:    strings.TrimSpace(stdout)
}

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
