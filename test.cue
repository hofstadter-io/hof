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
    @flow(test/hack)
    run: {
      @task(os.Exec)
      cmd: ["bash", "-c", scripts.DEV]
      script: "go test -v -cover pick_test.go"
      dir: "lib/structural"
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

