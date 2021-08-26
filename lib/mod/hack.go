package mod

import (
	"fmt"
	"io/ioutil"
	"strings"

	googithub "github.com/google/go-github/v38/github"
	"github.com/go-git/go-billy/v5/memfs"

	"github.com/hofstadter-io/hof/lib/mod/cache"
	"github.com/hofstadter-io/hof/lib/mod/parse/sumfile"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/lib/yagu/repos/github"
)

func Hack(lang string, args []string) error {
	fmt.Println("Hack", args)

	client, err := github.NewClient()
	if err != nil {
		return err
	}

	owner := args[0]
	repo := args[1]
	tag := args[2]

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
		fmt.Println(*t.Name, *t.Commit.SHA)
	}

	// Fetch and write to cache if tag found
	if T != nil {
		zReader, err := github.FetchTagZip(T)
		if err != nil {
			return fmt.Errorf("While fetching zipfile\n%w\n", err)
		}
		FS := memfs.New()

		err = yagu.BillyLoadFromZip(zReader, FS, true)
		if err != nil {
			return fmt.Errorf("While reading zipfile\n%w\n", err)
		}

		// fmt.Println("GOT HERE 1")

		err = cache.Write("hof", "github.com", owner, repo, tag, FS)
		if err != nil {
			return fmt.Errorf("While writing to cache\n%w\n", err)
		}

		// fmt.Println("GOT HERE 2")

		dirhash, err := yagu.BillyCalcHash(FS)
		if err != nil {
			return fmt.Errorf("While calculating dir hash\n%w\n", err)
		}

		modhash, err := yagu.BillyCalcFileHash("cue.mods", FS)
		if err != nil {
			return fmt.Errorf("While calculating mod hash\n%w\n", err)
		}

		S := sumfile.Sum{
			Mods: make(map[sumfile.Version][]string),
		}

		dver := sumfile.Version{
			Path: strings.Join([]string{"github.com", owner, repo}, "/"),
			Version: tag,
		}
		S.Add(dver, dirhash)

		mver := sumfile.Version{
			Path: strings.Join([]string{"github.com", owner, repo}, "/"),
			Version: strings.Join([]string{tag, "cue.mods"}, "/"),
		}
		S.Add(mver, modhash)

		fmt.Println("=====")
		out, err := S.Write()
		if err != nil {
			return err
		}
		fmt.Println(out)
		fmt.Println("=====")

		ioutil.WriteFile("cue.sums", []byte(out), 0644)

	}

	return nil
}
