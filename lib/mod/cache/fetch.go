package cache

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5/memfs"

	"github.com/hofstadter-io/hof/lib/yagu/repos/git"
	"github.com/hofstadter-io/hof/lib/yagu/repos/github"
	"github.com/hofstadter-io/hof/lib/yagu/repos/gitlab"
)

func Fetch(lang, mod, ver, pev string) (err error) {
	remote, owner, repo := parseModURL(mod)
	tag := ver

	dir := Outdir(lang, remote, owner, repo, tag)

	_, err = os.Lstat(dir)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return err
		}
		// not found, try fetching deps

		// determine private from ENV VAR in modder config (passed in as pev)
		private := MatchPrefixPatterns(os.Getenv(pev), mod)
		if err := fetch(lang, mod, ver, private); err != nil {
			return err
		}
	}

	// else we have it already
	return nil
}

func fetch(lang, mod, ver string, private bool) error {
	FS := memfs.New()

	remote, owner, repo := parseModURL(mod)
	tag := ver

	// TODO, how to deal with self-hosted / enterprise repos?

	switch remote {
	case "github.com":
		if err := github.Fetch(FS, owner, repo, tag, private); err != nil {
			return fmt.Errorf("While fetching from github\n%w\n", err)
		}

	case "gitlab.com":
		if err := gitlab.Fetch(FS, owner, repo, tag, private); err != nil {
			return fmt.Errorf("While fetching from gitlab\n%w\n", err)
		}

	default:
		if err := git.Fetch(FS, remote, owner, repo, tag, private); err != nil {
			return fmt.Errorf("While fetching from git\n%w\n", err)
		}
	}

	if err := Write(lang, remote, owner, repo, tag, FS); err != nil {
		return fmt.Errorf("While writing to cache\n%w\n", err)
	}

	return nil
}
