package hof

import "strings"

RepoRoot: { 
  @task(os.Exec)
  cmd: ["bash", "-c", "git rev-parse --show-toplevel"]
  stdout: string
  out: strings.TrimSpace(stdout)
}

watchAll: {
  @flow(watch/all)
  build: watchBuild
  test:  watchTest
}

watchBuild: {
  @flow(watch/build)

  RR: RepoRoot
  root: RR.out
  dirs: ["cmd","flow","lib","gen"]

  watch: {
    @task(fs.Watch)
    globs: [ for d in dirs { "\(root)/\(d)/**/*.go" } ]
    handler: {
      event?: _
      compile: {
        @task(os.Exec)
        cmd: ["go", "install", "\(root)/cmd/hof"]
      }
    }
  }
}
