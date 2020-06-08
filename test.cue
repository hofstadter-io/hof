package hof

//
////// Defined (partially) test configuration
//

#GoBaseTest: {
	skip: bool | *false

	sysenv: bool | *false
	env?: [string]: string
	...
}

#GoBashTest: #GoBaseTest & {
	script: string | *"""
	go test -cover ./
	"""
	...
}

#GoExecTest: #GoBaseTest & {
	command: string | *"go test -cover ./"
	...
}

//
////// Actual test configuration
//

// Test generated code
gen: _ @test(suite,gen)
gen: {
	// TODO before / after
	cmds: #GoBashTest @test(bash,gen/cmd)
	cmds: {
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
	self: #GoBashTest @test(bash,hsl)
	self: {
		dir: "script"
	}
}

hof: _ @test(suite,hof)
hof: {
	mod: #GoBashTest @test(bash,lib/mod)
	mod: {
		dir: "lib/mod"
	}
	st: #GoBashTest @test(bash,lib/st)
	st: {
		dir: "lib/structural"
	}
}
