package veg

flags: {
	// TODO, scope these since we are now injecting the entire package at the repo root
	repo: string | *"https://github.com/hofstadter-io/hof" @tag(repo)

	// todo, change this to "." when we move something to the index, if we ever really do?
	disk:   string | *"."             @tag(disk)
	src:    "repo" | *"disk" | string @tag(src,short=repo|disk)
	ref:    string | *"_next"         @tag(ref)
	adk:    string | *"../adk"        @tag(adk)
	dagger: string | *"../dagger"     @tag(dagger)

	goos: string @tag(goos,var=os)
	arch: string @tag(arch,var=arch)

	lsp: bool | *false

	ports: {
		gopls:  int | *4000 @tag(ports_gopls)
		cuepls: int | *4001 @tag(ports_cuepls)
	}

	registry: string | *"host.docker.internal:5000" @tag(reg)

	socket: string | *"unix:///var/run/docker.sock" @tag(socket)
	socket: "unix:///Users/tony/.colima/default/docker.sock"
}

gitFlags: {
	ci:          string @tag(ci,var=ci)
	gitRoot:     string @tag(gitRoot,var=gitRoot)
	gitCommit:   string @tag(gitCommit,var=gitCommit)
	gitShortSha: string @tag(gitShortSha,var=gitShortSha)
	gitBranch:   string @tag(gitBranch,var=gitBranch)
	gitTag:      string @tag(gitTag,var=gitTag)
	gitDirty:    string @tag(gitDirty,var=gitDirty)
}
