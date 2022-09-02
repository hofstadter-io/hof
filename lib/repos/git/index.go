package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"
)

// given a git url, and a req (required) version
// returns the minReference, allRefferences, and error
func IndexGitRemote(url, req string) (*plumbing.Reference, []*plumbing.Reference, error) {

	if !semver.IsValid(req) {
		return nil, nil, fmt.Errorf("Invalid SemVer v2 %q", req)
	}

	repo, err := NewRemote(url)
	if err != nil {
		return nil, nil, err
	}

	refs, err := repo.RemoteRefs()
	if err != nil {
		return nil, nil, err
	}

	var minRef *plumbing.Reference
	verS := ""

	// handle v0.0.0
	if req == "v0.0.0" {
		minRef, verS, err = findHead(refs)
		if err != nil {
			return nil, refs, err
		}
	} else {
		minRef, verS, err = findMin(req, refs)
		if err != nil {
			return nil, refs, err
		}
	}

	if verS == "" {
		return nil, refs, fmt.Errorf("Did not find compatible version for %s @ %s", url, req)
	}

	// fmt.Println("  found:", ref, hash)

	return minRef, refs, nil
}

func findHead(refs []*plumbing.Reference) (*plumbing.Reference, string, error) {
	// find HEAD ref
	headRef := ""
	for _, ref := range refs {
		fields := strings.Fields(ref.String())

		// HEAD is the only line with 3 fields
		if len(fields) < 3 {
			continue
		}

		v := fields[2]
		if v == "HEAD" {
			headRef = fields[1]
			break
		}
	}

	// find hash for HEAD
	for _, ref := range refs {
		fields := strings.Fields(ref.String())
		ver := fields[1]

		if ver == headRef {
			return ref, ver, nil
		}
	}

	return nil, "", nil
}

func findMin(req string, refs []*plumbing.Reference) (*plumbing.Reference, string, error) {

	var minR *plumbing.Reference
	min := ""
	for _, ref := range refs {
		fields := strings.Fields(ref.String())
		ver := fields[1]
		if strings.HasPrefix(ver, "refs/tags/") {
			ver = strings.TrimPrefix(ver, "refs/tags/")
			if ver[0:1] != "v" {
				ver = "v" + ver
			}

			if semver.IsValid(ver) {
				if semver.Compare(req, ver) <= 0 {
					// if this version is less than the current min, update
					if min == "" || semver.Compare(ver, min) < 0 {
						min = ver
						minR = ref
					}
				}
			}
		}
	}

	return minR, min, nil
}
