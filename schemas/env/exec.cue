@experiment(aliasv2)
package env

Exec: Step & {
	$kind: "exec"
	args: [...string]

	useEntrypoint?:  bool
	stdin?:          string
	redirectStdin?:  string
	redirectStdout?: string
	redirectStderr?: string
	expect?:         *"SUCCESS" | "FAILURE" | "ANY"

	experimentalPrivilegedNesting?: bool
	insecureRootCapabilities?:      bool
	expand?:                        bool
	noInit?:                        bool
}

Script: Exec & {
	script!: string
	args: ["sh", "-c", script]
}
Sh: Exec & {
	script!: string
	_script: """
		set -euo pipefail

		"""
	args: ["sh", "-c", _script + script]
}
Bash: Exec & {
	script!: string
	_script: """
		set -euo pipefail

		"""
	args: ["bash", "-c", _script + script]
}
Zsh: Exec & {
	script!: string
	_script: """
		set -euo pipefail

		"""
	args: ["zsh", "-c", _script + script]
}

// todo, think about how to hand stdio and redir to files,
// ideally they can be on the CUE types, but this is when we get into...
// the Fill CUE from Dagger results, continue eval'n CUE
// OG w/ Dagger, we did this through flow (?), avoid that here if possible
// just increasingly eval the value as much as we can?
// figure out what still needs to happen in Dagger, then do that
// this got more abstract than just Exec, applies to files/dir as well

Sync: Step & {
	$kind: "sync"
}

User: Step & {
	$kind: "user"
	name:  string
}

Workdir: Step & {
	$kind: "workdir"
	path:  string
}

Entrypoint: Step & {
	$kind: "entrypoint"
	args: [...string]

	keepDefaultArgs?: bool
}

DefaultArgs: Step & {
	$kind: "defaultArgs"
	args: [...string]
}

// sets the default terminal
DefaultTerm: Step & {
	$kind: "defaultTerm"
	args: [...string]

	experimentalPrivilegedNesting?: bool
	insecureRootCapabilities?:      bool
}

// starts an interactive terminal
Terminal: Step & {
	$kind: "terminal"
	args: [...string]
	experimentalPrivilegedNesting?: bool
	insecureRootCapabilities?:      bool
}

// TODO Vscode: Step & { ... } // or a family of them