package cache

import (
	"fmt"
	"strings"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"

	"github.com/hofstadter-io/hof/lib/repos/oci"
	"github.com/hofstadter-io/hof/lib/repos/remote"
)


func GetLatestTag(path string, pre bool) (string, error) {
	if debug {
		fmt.Println("cache.GetLatestTag", path, pre)
	}

	rmt, err := remote.Parse(path)
	if err != nil {
		return "", err
	}

	kind, err := rmt.Kind()
	if err != nil {
		return "", err
	}

	switch kind {
	case remote.KindGit:
		return GetLatestTagGit(path, pre)
	case remote.KindOCI:
		return GetLatestTagOCI(path, pre)
	}

	panic("cache.GetLatestTag: should not get here")
}

func GetLatestTagOCI(path string, pre bool) (string, error) {
	if debug {
		fmt.Println("cache.GetLatestTagOCI:", path, pre)
	}
	tags, err := oci.ListTags(path)
	if err != nil {
		return "", err
	}

	best := ""
	for _, n := range tags {
		// skip any tags which do not start with v
		if !strings.HasPrefix(n, "v") {
			continue
		}

		// maybe filter prereleases
		if !pre && semver.Prerelease(n) != "" {
			continue
		}

		// update best?
		if best == "" || semver.Compare(n, best) > 0 {
			best = n	
		}
	}

	return best, nil
}
	
func GetLatestTagGit(path string, pre bool) (string, error) {
	if debug {
		fmt.Println("cache.GetLatestTagGit:", path, pre)
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

// return the hash for the named ref, either git branch or OCI non-semver tag
func GetHashForNamedRef(path, ref string) (string, error) {
	if debug {
		fmt.Println("cache.GetHashForNamedRef", path, ref)
	}
	rmt, err := remote.Parse(path)
	if err != nil {
		return "", fmt.Errorf("remote parse: %w", err)
	}

	kind, err := rmt.Kind()
	if err != nil {
		return "", fmt.Errorf("remote kind: %w", err)
	}

	switch kind {
	case remote.KindGit:
		return GetBranchLatestHash(path, ref)
	case remote.KindOCI:
		return GetNamedTagHashOCI(path, ref)
	}

	panic("cache.GetLatestBranch: should not get here")
}

func GetNamedTagHashOCI(path, tag string) (string, error) {
	if debug {
		fmt.Println("cache.GetNamedTagHashOCI:", path, tag)
	}

	return oci.GetRefHash(path, tag)
}

// this function returns the commit hash for a branch
func GetBranchLatestHash(path, branch string) (string, error) {
	// sync
	_, err := FetchRepoSource(path, "")
	if err != nil {
		return branch, err
	}
	// open repo
	R, err := OpenRepoSource(path)
	if err != nil {
		return branch, fmt.Errorf("(gblh) open source error: %w for %s@%s", err, path, branch)
	}

	// get working tree
	wt, err := R.Worktree()
	if err != nil {
		return branch, fmt.Errorf("(gblh) worktree error: %w for %s@%s", err, path, branch)
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
	if debug {
		fmt.Println("cache.UpgradePsuedoVersion", path, ver)
	}

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
	nver, err := GetHashForNamedRef(path, ver)
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
