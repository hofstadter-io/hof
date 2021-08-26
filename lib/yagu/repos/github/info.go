package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v38/github"
)

func GetTagsSplit(client *github.Client, module string) ([]*github.RepositoryTag, error) {
	flds := strings.SplitN(module, "/", 1)
	domain, rest := flds[0], flds[1]

	if domain != "github.com" {
		return nil, fmt.Errorf("Github Tags Fetch called with non 'github.com' domain %q", module)
	}

	flds = strings.Split(rest, "/")
	owner, repo := flds[0], flds[1]
	tags, _, err := client.Repositories.ListTags(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func GetRepo(client *github.Client, owner, repo string) (*github.Repository, error) {
	r, _, err := client.Repositories.Get(context.Background(), owner, repo)
	return r, err
}

func GetBranch(client *github.Client, owner, repo, branch string) (*github.Branch, error) {
	b, _, err := client.Repositories.GetBranch(context.Background(), owner, repo, branch, true)
	return b, err
}

func GetBranches(client *github.Client, owner, repo, branch string) ([]*github.Branch, error) {
	bs, _, err := client.Repositories.ListBranches(context.Background(), owner, repo, nil)
	return bs, err
}

func GetTags(client *github.Client, owner, repo string) ([]*github.RepositoryTag, error) {
	tags, _, err := client.Repositories.ListTags(context.Background(), owner, repo, nil)
	return tags, err
}

