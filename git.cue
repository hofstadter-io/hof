package hof

import (
  "path"
  "strings"
)

RepoRoot: {
	@task(os.Exec)
	cmd: ["bash", "-c", "git rev-parse --show-toplevel"]
	stdout: string
	out:    strings.TrimSpace(stdout)
}

GitDiff: {
	@task(os.Exec)
  ref: string | *""
	cmd: ["bash", "-c", "git diff \(ref) --name-only"]
	stdout: string
	out:    strings.TrimSpace(stdout)
  files: strings.Split(out, "\n")
}

ShouldI: {
  globs: [...string]
  files: [...string]

  match: [
    for _, f in files
    for _, g in globs
    if path.Match(g, f, "unix") || strings.HasPrefix(f,g)
    { "\(f) -> \(g)" }
  ]

  yes: bool | *false
  if len(match) > 0 {
    yes: true
  }

}

shouldi_test: F={
  @flow(shouldi.test)

  globs: [...string] | *["*cue"]
  _g: string @tag(globs)
  if _g != _|_ {
    globs: strings.Split(_g, ",")
  }


  diff: GitDiff
  _shouldi: ShouldI & { files: diff.files }

  print: {
    @task(os.Stdout)
    shouldi: _shouldi & { globs: F.globs }

    if shouldi.yes {
      text: "yes\n"
    }
    if !shouldi.yes {
      text: "no\n"
    }
    @print()
  }
}
