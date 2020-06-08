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
	rm -rf .workdir
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
	cmds: #GoBashTest @test(bash,cmd)
	cmds: {
		dir: "cmd/hof/cmd"
	}
}

// Test Hof Linear Script (hls)
hls: _ @test(suite,hls)
hls: {
	self: #GoBashTest @test(bash,hls)
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
		skip: true
		dir: "lib/structural"
	}
}
