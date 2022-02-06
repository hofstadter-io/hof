// This flow gets an api code with OAuth workflow
package load

import (
  "encoding/json"

  "github.com/hofstadter-io/hof/examples/utils"
)

// twitch/auth/meta
meta: {
  @flow(meta,load) 

  vars: {
    RR: utils.RepoRoot
    root: RR.Out
    fn: "\(root)/flow/tests/repro/hof/data.json"
  }
  secrets: {
    env: { 
      FOO: _ @task(os.Getenv)
    } 
    foo: env.FOO
  }
}

// twitch/auth/load
thing: {
  @flow(thing,load)

  cfg: meta

  files: { 
    t: { filename: cfg.vars.fn } @task(os.ReadFile)
    j: json.Unmarshal(t.contents)
  } 
  data: files.j
  say: data.cow

  debug: { text: "load/data: " + files.t.contents} @task(os.Stdout)
}

