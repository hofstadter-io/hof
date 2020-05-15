package github

import (
	"context"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v30/github"
)

func NewClient() (client *github.Client, err error) {
	ctx := context.Background()

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)

	} else {
		client = github.NewClient(nil)
	}

	return client, err
}
