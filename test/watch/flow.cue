package hof

watchBuild: {
	@flow(watch/build)

	watch: {
		@task(os.Watch)
		globs: [ "./*.cue", "./template.*"]
		handler: {
			event?: _
			compile: {
				@task(os.Exec)
				cmd: ["bash", "-c", "hof gen data.cue -T 'template.txt;out.txt'"]
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
				text: "\(now.s) -> \(compile.cmd[2])\n"
			}
		}
	}
}

