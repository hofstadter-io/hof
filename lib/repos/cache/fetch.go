package cache

import (
	"fmt"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"

	"github.com/hofstadter-io/hof/lib/repos/git"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)

const hofPrivateVar = "HOF_PRIVATE"

func OpenRepoSource(path string) (*gogit.Repository, error) {
	if debug {
		fmt.Println("cache.OpenRepoSource:", path)
	}

	remote, owner, repo := utils.ParseModURL(path)
	dir := SourceOutdir(remote, owner, repo)
	return gogit.PlainOpen(dir)
}

func FetchRepoSource(path, ver string) (billy.Filesystem, error) {
	if debug {
		fmt.Println("cache.FetchRepoSource:", path)
	}

	remote, owner, repo := utils.ParseModURL(path)
	dir := SourceOutdir(remote, owner, repo)

	err := git.SyncSource(dir, remote, owner, repo, ver)
	if err != nil {
		return nil, err
	}

	FS := osfs.New(dir)

	return FS, nil
}

func GetLatestTag(path string, pre bool) (string, error) {
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

/*
func FindBestRef(url, ver string) (_,_,_ string, err error) {
	// fmt.Println("finding latest version for", url)
	if ver == "" {
		ver = "latest"
	}

	repo, err := git.NewRemote(url)
	if err != nil {
		return "", "", "", err
	}

	refs, err := repo.RemoteRefs()
	if err != nil {
		return "", "", "", err
	}

	refVal := ""
	refType := ""
	refCommit := ""

	for _, ref := range refs {
		parts := strings.Fields(ref.String())
		if len(parts) == 2 {
			commit, refStr := parts[0], parts[1]

			rs := strings.Split(refStr, "/")
			// HEAD & PRs, other?
			if len(rs) != 3 {
				continue
			}
			// ensure it is a ref
			if rs[0] != "refs" {
				continue
			}

			rType := rs[1]
			rVal := rs[2]

			// latest requires semver tags (optional leading 'v')
			if ver == "latest" {
				if rType != "tags" {
					continue
				}

				// add the optional 'v' if not present
				if rVal[0:1] != "v" {
					rVal = "v" + rVal
				}

				// continue if prerelease semver
				if semver.Prerelease(semver.Canonical(rVal)) != "" {
					continue
				}

				// update if empty
				if refVal == "" {
					refVal = rVal
					refCommit = commit
					refType = "tags"
				}

				if semver.IsValid(rVal) {
					if semver.Compare(rVal, refVal) > 0 {
						refVal = rVal
					}
				}

				// is this tag greater than a best so far? (semver wise)

			} else {
				// check if match
				//   break if found
				if rVal == ver {
					refVal = ver
					refCommit = commit
					refType = rType
					break
				}
			}
		}
	}

	// fmt.Println(refVal, refType, refCommit)

	return refVal, refType, refCommit, nil
}
*/
