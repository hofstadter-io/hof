package bbc

import (
	"fmt"
	"strings"

	"github.com/ktrysmt/go-bitbucket"
)

func GetTagsSplit(client *bitbucket.Client, module string) ([]bitbucket.RepositoryTag, error) {
	flds := strings.SplitN(module, "/", 1)
	domain, rest := flds[0], flds[1]

	if domain != "bitbucket.org" {
		return nil, fmt.Errorf("Bitbucket Tags Fetch called with non 'bitbucket.org' domain %q. Custom domains are a TODO", module)
	}

	flds = strings.Split(rest, "/")
	owner, repo := flds[0], flds[1]

	ro := &bitbucket.RepositoryTagOptions{
		Owner: owner,
		RepoSlug: repo,
	}
	tags, err := client.Repositories.Repository.ListTags(ro)
	if err != nil {
		return nil, err
	}
	return tags.Tags, nil
}

func GetRepo(client *bitbucket.Client, owner, repo string) (*bitbucket.Repository, error) {
	ro := &bitbucket.RepositoryOptions{
		Owner: owner,
		RepoSlug: repo,
	}
	r, err := client.Repositories.Repository.Get(ro)
	return r, err
}

func GetBranch(client *bitbucket.Client, owner, repo, branch string) (*bitbucket.RepositoryBranch, error) {
	ro := &bitbucket.RepositoryBranchOptions{
		Owner: owner,
		RepoSlug: repo,
	}
	b, err := client.Repositories.Repository.GetBranch(ro)
	return b, err
}

func GetBranches(client *bitbucket.Client, owner, repo, branch string) ([]bitbucket.RepositoryBranch, error) {
	ro := &bitbucket.RepositoryBranchOptions{
		Owner: owner,
		RepoSlug: repo,
	}
	bs, err := client.Repositories.Repository.ListBranches(ro)
	return bs.Branches, err
}

func GetTags(client *bitbucket.Client, owner, repo string) ([]bitbucket.RepositoryTag, error) {
	_, err := GetRepo(client, owner, repo)
	if err != nil {
		return nil, err
	}
	ro := &bitbucket.RepositoryTagOptions{
		Owner: owner,
		RepoSlug: repo,
	}
	tags, err := client.Repositories.Repository.ListTags(ro)
	return tags.Tags, err
}

