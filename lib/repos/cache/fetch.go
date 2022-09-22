package cache

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"golang.org/x/mod/semver"

	"github.com/hofstadter-io/hof/lib/repos/bbc"
	"github.com/hofstadter-io/hof/lib/repos/git"
	"github.com/hofstadter-io/hof/lib/repos/github"
	"github.com/hofstadter-io/hof/lib/repos/gitlab"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)

const hofPrivateVar = "HOF_PRIVATE"

func FetchRepoToMem(url, ver string) (billy.Filesystem, error) {
	fmt.Println("downloading: ", url, ver)
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

