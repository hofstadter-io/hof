package hof

// Test generated code
gen: _ @test(suite,gen)
gen: {
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

// Test Hof Linear Script (hls)
hls: _ @test(suite,hls)
hls: {
	self: _ @test(script,hsl,self)
	self: {
		dir: "script"
		scripts: [ "testdata/*.txt" ]
	}
}

hof: _ @test(suite,hof)
hof: {
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
