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
	scdir: "\(back)sonar-reports/go/\(dir)"
	script: string | *"""
	rm -rf .workdir
	mkdir -p \(scdir)
	go test -cover ./ -coverprofile \(scdir)/cover.out -json > \(scdir)/tests.out
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
	cmds: #GoBashTest @test(bash,cmd,test)
	cmds: {
		dir: "cmd/hof/cmd"
	}
	cmdsC: #GoBashCover @test(bash,cmd,cover)
	cmdsC: {
		dir: "cmd/hof/cmd"
	}
}

// Test Hof Linear Script (hls)
hls: _ @test(suite,script,hls)
hls: {
	self: #GoBashTest @test(bash,test)
	self: {
		dir: "script"
	}
	selfC: #GoBashCover @test(bash,cover)
	selfC: {
		dir: "script"
	}
}

// Test Hof Linear Script (hls)
hsh: _ @test(suite,shell,hsh,sh)
hsh: {
	self: #GoBashTest @test(bash,test)
	self: {
		dir: "shell"
	}
	selfC: #GoBashCover @test(bash,cover)
	selfC: {
		dir: "shell"
	}
}

hof: _ @test(suite,lib,hof)
hof: {

	mod: #GoBashTest @test(bash,lib/mod,test)
	mod: {
		dir: "lib/mod"
	}
	modC: #GoBashCover @test(bash,lib/mod,cover)
	modC: {
		dir: "lib/mod"
	}

	st: #GoBashTest @test(bash,lib/st,test)
	st: {
		skip: true
		dir: "lib/structural"
	}
	stC: #GoBashCover @test(bash,lib/st,cover)
	stC: {
		skip: true
		dir: "lib/structural"
	}

}
