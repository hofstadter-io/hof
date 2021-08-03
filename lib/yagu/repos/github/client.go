package github

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v30/github"
)

const TokenEnv = "GITHUB_TOKEN"

func NewClient() (client *github.Client, err error) {
	ctx := context.Background()

	if token := os.Getenv(TokenEnv); token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)

	} else {
		fmt.Println("Got Here GitHub")
		client = github.NewClient(nil)
	}

	return client, err
}
