@experiment(aliasv2)
package inception

import (
  "github.com/hofstadter-io/hof/catalogs/env/bases"
  "github.com/hofstadter-io/hof/catalogs/env/packs"
  "github.com/hofstadter-io/hof/catalogs/env/steps"
  // "github.com/hofstadter-io/hof/examples/env/veg"
  "github.com/hofstadter-io/hof/schemas/env"
)

let root = self

// _veg:   veg   & { flags: root.flags}
// these two aren't needed yet, but that's the plan
_packs: packs & { flags: root.flags}
_steps: steps & { flags: root.flags}


turtles: {
  // any container
  dev: env.#Container & {
    @env(turtles-dev)
    from: bases.debian13.default

    steps: [
      // use zsh
      _steps.tool.zsh.customize,

      // add hof/veg
      // _veg.hof.File["linux-arm64"],

      // add a bunch of tools
      _packs.containers.docker.cli.install,
      _steps.tool.dagger.cli,
			_steps.tool.k8s.kubectl,
			_steps.tool.k8s.helm,
			_steps.tool.k8s.crane,
			_steps.tool.k8s.kind,

      // add the socket for inception
      env.UnixSocket & { path: "/var/run/docker.sock", source: turtles.socket },
    ]
  }

  // get a socket from the host, inception will be painfully slow otherwise
  #socket: string | *"unix:///var/run/docker.sock"
  socket: env.#HostSocket & {
    @env(turtles-sock)
    name: "turtles-sock"
    path: #socket
  }

}

flags: {
	// TODO, scope these since we are now injecting the entire package at the repo root
	repo: string | *"https://github.com/hofstadter-io/hof" @tag(repo)

	// todo, change this to "." when we move something to the index, if we ever really do?
	disk:  string | *"."              @tag(disk,var=gitRoot)
	src:    "repo" | *"disk" | string @tag(src,short=repo|disk)
	ref:    string | *"_next"          @tag(ref)
	adk:    string | *"../adk"         @tag(adk)
	dagger: string | *"../dagger"      @tag(dagger)

	goos: string @tag(goos,var=os)
	arch: string @tag(arch,var=arch)

	lsp: bool | *false

	ports: {
		gopls:  int | *4000 @tag(ports_gopls)
		cuepls: int | *4001 @tag(ports_cuepls)
	}

	registry: string | *"host.docker.internal:5000" @tag(reg)
}

gitFlags: {
	ci: string @tag(ci,var=ci)
	gitRoot: string @tag(gitRoot,var=gitRoot)
	gitCommit: string @tag(gitCommit,var=gitCommit)
	gitShortSha: string @tag(gitShortSha,var=gitShortSha)
	gitBranch: string @tag(gitBranch,var=gitBranch)
	gitTag: string @tag(gitTag,var=gitTag)
	gitDirty: string @tag(gitDirty,var=gitDirty)
}
