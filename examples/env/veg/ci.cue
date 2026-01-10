@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

ci: {

  // These are set to developer local by default
  // in CI, these are different git repo checkouts
  // - for a PR, the target and current
  // - for a branch w/o PR, main branch, until bfg
  #prev: _ | *src.repo
  #next: _ | *src.disk

  // this is what we run to verify the cue mod
  cuemod: env.#Shouldi & {
    @env(ci-cuemod-shouldi)
    changes: env.#Changes & {
      prev: #prev
      next: #next
    }
    include: root.dist.cuemod.include
    then: env.#Task & {
      @env(ci-cuemod-task)
      steps: [
        // fmt/lint
        [_tester & {#cmd: "cue fmt ./..."}],

        // test
        [root.cmd.test.tasks.go, root.cmd.test.tasks.goveti]
      ]
    }
  }

  cuefmt: {
    fmtd: env.#Container & {
      @env(ci-cuemod-fmtd)
      from: root.ctr.dev
      steps: [
        env.Dir & {source: root.src.cuemod},
        env.Sh & {script: "cue fmt ./..."},
      ]
    }
  }

  _tester: env.#Container & {
    #cmd: string
    from: env.#Container & {
      @id(ci-tester-with-src)
      from: root.ctr.dev
      steps: [
        
      ]
    }
    steps: [
      env.Bash & {script: "\(#cmd)"},
    ]
  }

}