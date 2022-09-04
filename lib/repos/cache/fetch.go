package cache

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"

	"github.com/hofstadter-io/hof/lib/repos/bbc"
	"github.com/hofstadter-io/hof/lib/repos/git"
	"github.com/hofstadter-io/hof/lib/repos/github"
	"github.com/hofstadter-io/hof/lib/repos/gitlab"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)

const hofPrivateVar = "HOF_PRIVATE"

func FetchRepoToMem(url, ver string) (billy.Filesystem, error) {
	fmt.Println("fetch: ", url, ver)
	FS := memfs.New()

	remote, owner, repo := utils.ParseModURL(url)

	// check if in cache?

	// TODO, how to deal with self-hosted / enterprise repos?
	private := utils.MatchPrefixPatterns(os.Getenv(hofPrivateVar), url)

	// TODO, make an interface for repo hosts
	switch remote {
	case "github.com":
		if err := github.Fetch(FS, owner, repo, ver, private); err != nil {
			return FS, fmt.Errorf("While fetching from github\n%w\n", err)
		}

	case "gitlab.com":
		if err := gitlab.Fetch(FS, owner, repo, ver, private); err != nil {
			return FS, fmt.Errorf("While fetching from gitlab\n%w\n", err)
		}

	case "bitbucket.org":
		if err := bbc.Fetch(FS, owner, repo, ver, private); err != nil {
			return FS, fmt.Errorf("While fetching from bitbucket\n%w\n", err)
		}

	default:
		if err := git.Fetch(FS, remote, owner, repo, ver, private); err != nil {
			return FS, fmt.Errorf("While fetching from git\n%w\n", err)
		}
	}

	return FS, nil
}
