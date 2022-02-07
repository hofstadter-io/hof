package hof

GoTest: {
  @task(os.Exec)
  cmd: ["bash", "-c", scripts.DEV]
}

scripts: {
  DEV: string | *"""
  rm -rf .workdir
  go test -cover ./
  """

  CI: string | *"""
  rm -rf .workdir
  go test -cover ./ -coverprofile cover.out -json > tests.json
  """
}

watchTest: {
  @flow(watch/test)

  watch: {
    @task(fs.Watch)
    first: true
    globs: [
      "lib/structural/**/*.*",
    ]
    handler: {
      event?: _
      compile: {
        @task(os.Exec)
        cmd: ["hof", "flow", "-f", "test/hack"]
      }
    }
  }
}

tests: {
  // want to discover nested too
  // @flow(test)

  hack: {
    test: string | *"TestMainFlow" @tag(test)
    @flow(test/hack)
    prt: { text: "testing: \(test)\n" } @task(os.Stdout)
    run: {
      @task(os.Exec)
      cmd: ["bash", "-c", script]
      dir: "flow"
      script: """
      rm -rf .workdir
      go test -run \(test) . 
      """
    }
  }

  flow: {
    @flow(test/flow)
    run: GoTest & {
      // dir: "lib/flow" // panics, segfault
      dir: "flow"
    }
  }

  st: {
    @flow(test/st)
    run: GoTest & {
      dir: "lib/structural"
    }
  }

  mods: {
    @flow(test/mods)
    run: GoTest & {
      dir: "lib/mods"
    }
  }
}

