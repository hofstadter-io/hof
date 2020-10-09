package cli

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#CliReleases: schema.#GoReleaser & {
	Disable: false
	Draft:    false
	Author:   "Hofstadter, Inc"
	Homepage: "https://hofstadter.io"

	GitHub: {
		Owner: "hofstadter-io"
		Repo:  "hof"
	}

	Docker: {
		Maintainer: "Hofstadter, Inc <open-source@hofstadter.io>"
		Repo:       "hofstadter"
	}
}
