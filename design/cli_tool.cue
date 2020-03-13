package cli

import (
	"path"
  "tool/exec"
  "tool/file"
)

command: gen: {

  var: {
    outdir: "output/"
  }

  for i, F in GEN.Out {

    task: "mkdir-\(i)": exec.Run & {
      cmd:    ["mkdir", "-p", var.outdir + path.Dir(F.Filename)]
      stdout: string
    }

    task: "write-\(i)": file.Create & {
      filename: var.outdir + F.Filename
      contents: F.Out
      deps: [ task["mkdir-\(i)"].stdout ]
    }

  }

}

