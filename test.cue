package hof

#GoTest: {
	script: """
	go test -cover ./
	"""

	skip: bool | *false
	...
}

// Test generated code
gen: _ @test(suite,gen)
gen: {
	// TODO before / after
	cmds: _ @test(bash,gen/cmd)
	cmds: #GoTest & {
		dir: "cmd/hof/cmd"
		env: {
			HELLO: "WORLD"
			FOO: "$FOO"
		}
	}
}

// Test Hof Linear Script (hls)
hls: _ @test(suite,hls)
hls: {
	self: #GoTest @test(bash,hsl)
	self: {
		dir: "script"
	}
}

hof: _ @test(suite,hof)
hof: {
	mod: #GoTest @test(bash,lib/mod)
	mod: {
		dir: "lib/mod"
	}
	st: #GoTest @test(bash,lib/st)
	st: {
		dir: "lib/structural"
	}
}
