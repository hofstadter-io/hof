package fetch

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

func FetchRepoToTmp(url, ver string, private bool) (err error) {
	//FS, err := FetchRepoToMem(url, ver)
	//if err != nil {
		//return err
	//}

	// write to temp dir
	//if err := Write(lang, remote, owner, repo, ver, FS); err != nil {
		//return fmt.Errorf("While writing to cache\n%w\n", err)
	//}

	// else we have it already
	return nil
}

func FetchRepoToMem(url, ver string) (billy.Filesystem, error) {
	fmt.Println("fetch: ", url, ver)
	FS := memfs.New()

	remote, owner, repo := utils.ParseModURL(url)

	// TODO, how to deal with self-hosted / enterprise repos?
	private := utils.MatchPrefixPatterns(os.Getenv(hofPrivateVar), url)

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
