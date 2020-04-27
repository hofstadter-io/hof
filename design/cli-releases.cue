package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#CliReleases: schema.#GoReleaser & {
  Disabled: false
  Draft: false
  Author:   "Tony Worm"
  Homepage: "https://github.com/hofstadter-io/hof"

  GitHub: {
    Owner: "hofstadter-io"
    Repo:  "hof"
  }

  Docker: {
    Maintainer: "Hofstadter, Inc <open-source@hofstadter.io>"
    Repo: "hofstadter"
  }
}

