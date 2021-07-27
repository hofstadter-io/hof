package gitlab

import (
	"archive/zip"
	"bytes"
	"fmt"

	"github.com/xanzy/go-gitlab"
)

func fetchShaZip(client *gitlab.Client, pid interface{}, sha string) (*zip.Reader, error) {
	format := "zip"
	data, _, err := client.Repositories.Archive(pid, &gitlab.ArchiveOptions{
		Format: &format,
		SHA:    &sha,
	})
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)

	return zip.NewReader(r, int64(len(data)))
}

func FetchZip(client *gitlab.Client, owner, repo, tag string) (*zip.Reader, error) {
	pid := fmt.Sprintf("%s/%s", owner, repo)

	var sha string

	if tag == "v0.0.0" {
		bs, _, err := client.Branches.ListBranches(pid, nil)
		if err != nil {
			return nil, err
		}

		var branch *gitlab.Branch

		for _, candidate := range bs {
			if candidate.Default {
				branch = candidate

				break
			}
		}

		if branch == nil {
			return nil, fmt.Errorf("Could not find default branch for repository %s", pid)
		}

		sha = branch.Commit.ID
	} else {
		t, _, err := client.Tags.GetTag(pid, tag)
		if err != nil {
			return nil, err
		}

		sha = t.Commit.ID
	}

	return fetchShaZip(client, pid, sha)
}
