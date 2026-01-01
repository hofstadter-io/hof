package veg

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)

src: {
	repo: env.#GitRepo & {
		@env()
		name: string | *"repo"
		url:  flags.repo
		ref:  flags.ref
	}
	local: env.#HostDir & {
		@env()
		name: string | *"local"
		path: flags.local
	}
	adk: env.#HostDir & {@env(), path: flags.adk}
	dagger: env.#HostDir & {@env(), path: flags.dagger}

	// setup code base on flags and value
	code: {@env(), name: "code"}
	if flags.src == "repo" {code: repo}
	if flags.src == "local" {code: local}
	if strings.HasPrefix(flags.src, "https://") {
		code: env.#GitRepo & {url: flags.src}
	}
	if strings.HasPrefix(flags.src, ".") {
		code: env.#HostDir & {path: flags.src}
	}

	extn: {
		vscode: {}
	}
}
