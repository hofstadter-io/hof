package cache

import (
	"fmt"
	"strings"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"
)


func GetLatestTag(path string, pre bool) (string, error) {
	if debug {
		fmt.Println("cache.GetLatestBranch", path, pre)
	}
	_, err := FetchRepoSource(path, "")
	if err != nil {
		return "", err
	}

	R, err := OpenRepoSource(path)
	if err != nil {
		return "", err
	}

	refs, err := R.Tags()
	if err != nil {
		return "", err
	}

	best := ""
	refs.ForEach(func (ref *plumbing.Reference) error {
		n := ref.Name().Short()

		// skip any tags which do not start with v
		if !strings.HasPrefix(n, "v") {
			return nil
		}

		// maybe filter prereleases
		if !pre && semver.Prerelease(n) != "" {
			return nil
		}

		// update best?
		if best == "" || semver.Compare(n, best) > 0 {
			best = n	
		}
		return nil
	})
	
	return best, nil
}

func GetLatestBranch(path, branch string) (string, error) {
	if debug {
		fmt.Println("cache.GetLatestBranch", path, branch)
	}
	// sync
	_, err := FetchRepoSource(path, "")
	if err != nil {
		return branch, err
	}
	// open repo
	R, err := OpenRepoSource(path)
	if err != nil {
		return branch, fmt.Errorf("(glb) open source error: %w for %s@%s", err, path, branch)
	}

	// get working tree
	wt, err := R.Worktree()
	if err != nil {
		return branch, fmt.Errorf("(glb) worktree error: %w for %s@%s", err, path, branch)
	}

	// try to checkout branch
	err = wt.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.NewRemoteReferenceName("origin", branch),
		Force: true,
	})
	if err == nil {
		h, err := R.Head()
		if err != nil {
			return branch, err
		}
		return h.Hash().String(), nil
	}

	return branch, nil
}

func UpgradePseudoVersion(path, ver string) (s string, err error) {
	// semver tag?
	if semver.IsValid(ver) {
		return ver, nil
	}

	if ver == "latest" || ver == "next" {
		ver, err = GetLatestTag(path, ver == "next")
		if err != nil {
			return ver, err
		}
	}

	// branch? need to find commit
	nver, err := GetLatestBranch(path, ver)
	if err != nil {
		return ver, err
	}
	if nver != "" {
		ver = nver
	}

	if !strings.HasPrefix(ver, "v") {
		now := time.Now().UTC().Format("20060102150405")
		ver = fmt.Sprintf("v0.0.0-%s-%s", now, ver)
	}


	return ver, nil
}

