package cache

import (
	"fmt"
	"os"
	"strings"

	googithub "github.com/google/go-github/v30/github"
	"github.com/go-git/go-billy/v5/memfs"

	"github.com/hofstadter-io/hof/lib/mod/repos/github"
	"github.com/hofstadter-io/hof/lib/mod/util"
)

func Fetch(lang, mod, ver string) (err error) {
	flds := strings.Split(mod, "/")
	remote := flds[0]
	owner := flds[1]
	repo := flds[2]
	tag := ver

	dir := Outdir(lang, remote, owner, repo, tag)

	_, err = os.Lstat(dir)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return err
		}
		// not found
		fetch(lang, mod, ver)
	}

	// else we have it already
	// fmt.Println("Found in cache")

	return nil
}

func fetch(lang, mod, ver string) error {
	flds := strings.Split(mod, "/")
	remote := flds[0]
	owner := flds[1]
	repo := flds[2]
	tag := ver

	switch remote {
	case "github.com":
		return fetchGitHub(lang, owner, repo, tag)

	default:
		return fmt.Errorf("Unknown remote: %q in %s", remote, mod)
	}
}

func fetchGitHub(lang, owner, repo, tag string) error {
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
			fmt.Printf("FOUND  ")
		}
		// fmt.Println(*t.Name, *t.Commit.SHA)
	}

	if T != nil {
		zReader, err := github.FetchTagZip(client, T)
		if err != nil {
			return fmt.Errorf("While fetching zipfile\n%w\n", err)
		}
		FS := memfs.New()

		err = util.BillyLoadFromZip(zReader, FS, true)
		if err != nil {
			return fmt.Errorf("While reading zipfile\n%w\n", err)
		}

		// fmt.Println("GOT HERE 1")

		err = Write(lang, "github.com", owner, repo, tag, FS)
		if err != nil {
			return fmt.Errorf("While writing to cache\n%w\n", err)
		}
	}

	return nil
}
