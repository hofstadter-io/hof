package hof

import "strings"

//
////// Defined (partially) test configuration
//

#GoBaseTest: {
	skip: bool | *false

	sysenv: bool | *false
	env?: [string]: string

	dir: string
	...
}

#GoBashTest: #GoBaseTest & {
	dir: string
	script: string | *"""
	rm -rf .workdir
	go test -cover ./
	"""
	...
}

#GoBashCover: #GoBaseTest & {
	dir: string
	back: strings.Repeat("../", strings.Count(dir, "/") + 1)
	script: string | *"""
	rm -rf .workdir
	go test -cover ./ -coverprofile cover.out -json > tests.json
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
	cmds: #GoBashTest @test(bash,test,cmd)
	cmds: {
		dir: "cmd/hof/cmd"
	}
	cmdsC: #GoBashCover @test(bash,cover,cmd)
	cmdsC: {
		dir: "cmd/hof/cmd"
	}
}

// Test Hof Linear Script (hls)
hls: _ @test(suite,hls)
hls: {
	runtime: #GoBashTest @test(bash,test,runtime)
	runtime: {
		dir: "script/runtime"
	}
	runtimeC: #GoBashCover @test(bash,cover,runtime)
	runtimeC: {
		dir: "script/runtime"
	}

	shell: #GoBashTest @test(bash,test,shell)
	shell: {
		dir: "script/shell"
	}
	shellC: #GoBashCover @test(bash,cover,shell)
	shellC: {
		dir: "script/shell"
	}

	script: #GoBashTest @test(bash,test,script)
	script: {
		dir: "script"
	}
	scriptC: #GoBashCover @test(bash,cover,script)
	scriptC: {
		dir: "script"
	}
}

lib: _ @test(suite,lib)
lib: {

	mod: #GoBashTest @test(bash,test,mod)
	mod: {
		dir: "lib/mod"
	}
	modC: #GoBashCover @test(bash,cover,mod)
	modC: {
		dir: "lib/mod"
	}

	st: #GoBashTest @test(bash,test,st)
	st: {
		dir: "lib/structural"
	}
	stC: #GoBashCover @test(bash,cover,st)
	stC: {
		dir: "lib/structural"
	}

}
