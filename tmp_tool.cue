package cli

import (
	"path"
	"tool/cli"
	"tool/exec"
	"tool/file"

  "github.com/hofstadter-io/cuelib/template"
)

command: gen: {
	task: step_1: exec.Run & {
		cmd: ["cue", "render"]
		stdout: string
	}
	task: step_2: exec.Run & {
		$after: task.step_1
		cmd: ["cue", "format"]
		stdout: string
	}
}

command: render: {

	var: {
		outdir: Outdir
	}

	for i, F in GEN.Out {

    if F.Filename != _|_ {
      TMP = {
        if F.Alt == _|_ {
          Out: (template.RenderTemplate & { Template: F.Template, Values: F.In}).Out
        }
        if F.Alt != _|_ {
          Out: (template.AltDelimTemplate & { Template: F.Template, Values: F.In}).Out
        }
      }

      task: "mkdir-\(i)": exec.Run & {
        cmd: ["mkdir", "-p", var.outdir + path.Dir(F.Filename)]
        stdout: string
      }

      task: "write-\(i)": file.Create & {
        deps: [ task["mkdir-\(i)"].stdout]

        filename: var.outdir + F.Filename
        contents: TMP.Out
        stdout:   string
      }

      task: "print-\(i)": cli.Print & {
        deps: [ task["write-\(i)"].stdout]
        text: task["write-\(i)"].filename
      }
    } 

	}

}

command: format: {
	var: {
		outdir: "."
	}

	task: shell: exec.Run & {
		cmd: ["bash", "-c", "cd \(var.outdir) && goimports -w -l ."]
		stdout: string
	}
}

command: build: {
	var: {
		outdir: Outdir
	}

	task: shell: exec.Run & {
		cmd: ["bash", "-c", "cd \(var.outdir) && go build -o \(CLI.Name) ."]
		stdout: string
	}

}
