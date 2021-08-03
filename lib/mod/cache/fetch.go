package cache

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	googithub "github.com/google/go-github/v30/github"

	"github.com/hofstadter-io/hof/lib/yagu"
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
		private := MatchPrefixPatterns(os.Getenv(pev), mod)
		if err := fetch(lang, mod, ver, private); err != nil {
			return err
		}
	}

	// else we have it already
	// fmt.Println("Found in cache")

	return nil
}

func fetch(lang, mod, ver string, private bool) error {
	remote, owner, repo := parseModURL(mod)
	tag := ver

	if private {
		return fetchGit(lang, remote, owner, repo, tag, true)
	} else if remote == "github.com" {
		return fetchGitHub(lang, owner, repo, tag)
	} else if remote == "gitlab.com" {
		return fetchGitLab(lang, owner, repo, tag)
	}
	return fetchGit(lang, remote, owner, repo, tag, private)
}

func fetchGit(lang, remote, owner, repo, tag string, private bool) error {
	FS := memfs.New()

	if err := git.FetchGit(FS, remote, owner, repo, tag, private); err != nil {
		return fmt.Errorf("While fetching from git\n%w\n", err)
	}

	if err := Write(lang, remote, owner, repo, tag, FS); err != nil {
		return fmt.Errorf("While writing to cache\n%w\n", err)
	}

	return nil
}

func fetchGitLab(lang, owner, repo, tag string) (err error) {
	FS := memfs.New()
	client, err := gitlab.NewClient()
	if err != nil {
		return err
	}

	zReader, err := gitlab.FetchZip(client, owner, repo, tag)
	if err != nil {
		return fmt.Errorf("While fetching from GitLab\n%w\n", err)
	}

	if err := yagu.BillyLoadFromZip(zReader, FS, true); err != nil {
		return fmt.Errorf("While reading zipfile\n%w\n", err)
	}

	if err := Write(lang, "gitlab.com", owner, repo, tag, FS); err != nil {
		return fmt.Errorf("While writing to cache\n%w\n", err)
	}

	return nil
}

func fetchGitHub(lang, owner, repo, tag string) (err error) {
	FS := memfs.New()

	if tag == "v0.0.0" {
		err = fetchGitHubBranch(FS, lang, owner, repo, "")
	} else {
		err = fetchGitHubTag(FS, lang, owner, repo, tag)
	}
	if err != nil {
		return fmt.Errorf("While fetching from github\n%w\n", err)
	}

	err = Write(lang, "github.com", owner, repo, tag, FS)
	if err != nil {
		return fmt.Errorf("While writing to cache\n%w\n", err)
	}

	return nil
}
func fetchGitHubBranch(FS billy.Filesystem, lang, owner, repo, branch string) error {
	client, err := github.NewClient()
	if err != nil {
		return err
	}

	if branch == "" {
		r, err := github.GetRepo(client, owner, repo)
		if err != nil {
			return err
		}
		branch = *r.DefaultBranch

		fmt.Printf("%#+v\n", *r)
	}

	// fmt.Println("Fetch github BRANCH", lang, owner, repo, branch)

	zReader, err := github.FetchBranchZip(client, owner, repo, branch)
	if err != nil {
		return fmt.Errorf("While fetching branch zipfile for %s/%s@%s\n%w\n", owner, repo, branch, err)
	}

	err = yagu.BillyLoadFromZip(zReader, FS, true)
	if err != nil {
		return fmt.Errorf("While reading branch zipfile\n%w\n", err)
	}

	return nil
}
func fetchGitHubTag(FS billy.Filesystem, lang, owner, repo, tag string) error {
	// fmt.Println("Fetch github TAG", lang, owner, repo, tag)
	client, err := github.NewClient()
	if err != nil {
		return err
	}

	tags, err := github.GetTags(client, owner, repo)
	if err != nil {
		return err
	}

	// The tag we are looking for
	var T *googithub.RepositoryTag
	for _, t := range tags {
		if tag != "" && tag == *t.Name {
			T = t
			// fmt.Printf("FOUND  %v\n", *t.Name)
		}
	}
	if T == nil {
		return fmt.Errorf("Did not find tag %q for 'https://github.com/%s/%s' @%s", tag, owner, repo, tag)
	}

	zReader, err := github.FetchTagZip(client, T)
	if err != nil {
		return fmt.Errorf("While fetching tag zipfile\n%w\n", err)
	}

	err = yagu.BillyLoadFromZip(zReader, FS, true)
	if err != nil {
		return fmt.Errorf("While reading tag zipfile\n%w\n", err)
	}

	return nil
}
