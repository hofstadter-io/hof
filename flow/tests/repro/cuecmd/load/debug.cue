// This flow gets an api code with OAuth workflow
package load

import (
  "encoding/json"

  "hof.io/repro/utils"
)

meta: {
  vars: {
    RR: utils.RepoRoot
    root: RR.Out
    fn: "\(root)/examples/repro/cuecmd/data.json"
  }
}

load: {
  cfg: meta

  read: {
    $id: "tool/file.Read"
    filename: cfg.vars.fn
    contents: string
  }

  data: json.Unmarshal(read.contents)
  say: data.cow

  print: {
    $id: "tool/cli.Print"
    text: read.contents
  }
}
