package cli

import (
	"path"
  "tool/cli"
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
      deps: [ task["mkdir-\(i)"].stdout ]

      filename: var.outdir + F.Filename
      contents: F.Out
      stdout: string
    }

    task: "print-\(i)": cli.Print & {
      deps: [ task["write-\(i)"].stdout ]
      text: task["write-\(i)"].filename
    }

  }

  task: format: exec.Run & {
    cnt : len(GEN.Out) - 1
    deps: [ 
      task["write-0"].stdout,
      task["write-\(cnt)"].stdout
    ]
    cmd: ["bash", "-c", "cd \(var.outdir) && go fmt ./..."]
    stdout: string
  }

  task: print: cli.Print & {
    text: task.format.stdout
  }
}

