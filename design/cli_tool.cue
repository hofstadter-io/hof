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

  task: clean: exec.Run & {
    cmd:    ["rm", "-rf", var.outdir]
    stdout: string  // capture stdout
  }
  task: mkdir: exec.Run & {
    cmd:    ["mkdir", "-p", var.outdir]
    stdout: string
    deps: [ task["clean"].stdout ]
  }

  for i, FS in GEN._Out {
		for j, F in FS {
			task: "mkdir-\(i)-\(j)": exec.Run & {
				cmd:    ["mkdir", "-p", var.outdir + path.Dir(F._Filename)]
				stdout: string
				deps: [ task["mkdir"].stdout ]
			}
			task: "write-\(i)-\(j)": file.Create & {
				filename: var.outdir + F._Filename
				contents: F._Out
				deps: [ task["mkdir-\(i)-\(j)"].stdout ]
			}

		}
  }

}

