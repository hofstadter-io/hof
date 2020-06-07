package hof

generated: _ @test(suite,gen/cli)
generated: {
	// TODO before / after
	cmds: _ @test(script)
	cmds: {
		dir: "cmd/hof/cmd"
		scripts: [ "**/*.txt" ]
		env: {
			HELLO: "WORLD"
			FOO: "$FOO"
		}
	}
}

human: _ @test(suite,ppl)
human: {
	mod: _ @test(script,lib/mod)
	mod: {
		dir: "lib/mod"
		scripts: [ "testdata/*.txt" ]
	}
	st: _ @test(exec=bash,lib/st)
	st: {
		dir: "lib/structural"
		script: """
		go test -cover ./
		"""
	}
}
