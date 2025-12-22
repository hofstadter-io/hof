package env


import (
	"github.com/hofstadter-io/hof/schemas"
	"github.com/hofstadter-io/hof/schemas/common"
)

// Definition for a container env
Container: {
	schemas.Hof
	#hof: env: {
        root: true
        kind: "container"
    }
	Name: common.NameLabel

    From: string
    Labels: [string]: string
    Steps: [...Step]
}

StepKinds: [
    // running inside
    "exec",

    // filesys related
    "file",
    "files",
    "dir",

    // envs & secrets
    "env",
    "envfile",
    "secret",

    // external resources
    "mount", // consolidated, we may want to split them?

    // runtime stuff
    "expose",
    "entrypoint",
    "args",
    "term", // default dagger term setup
]

Step: {
    kind: or(StepKinds)
}

StepExec: {
    kind: "exec"
}

StepFile: {
    kind: "file"
}

StepFiles: {
    kind: "files"
}

StepDir: {
    kind: "dir"
}

StepEnv: {
    kind: "env"
}

StepEnvfile: {
    kind: "envfile"
}

StepSecret: {
    kind: "secret"
}

StepMount: {
    kind: "mount"
    // cache, dir, file, secret, temp, service (?)
}

StepExpose: {
    kind: "expose"
}

StepEntrypoint: {
    kind: "entrypoint"
}

StepArgs: {
    kind: "args"
}

StepTerm: {
    kind: "term"
}

