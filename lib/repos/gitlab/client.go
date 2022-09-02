package gitlab

import (
	"fmt"
	"os"

	"github.com/xanzy/go-gitlab"
)

const TokenEnv = "GITLAB_TOKEN"

func NewClient(private bool) (client *gitlab.Client, err error) {
	// TODO, there are multiple NewClient<> methods
	// how to determine which to use and the inputs to them
	// https://pkg.go.dev/github.com/xanzy/go-gitlab#NewOAuthClient
	// Noting also that we prefer auth over non-auth for API rate limits

	token := os.Getenv(TokenEnv)
	if private && token == "" {
		return nil, fmt.Errorf("Private module requested and no auth token available for %s", TokenEnv)
	}
	return gitlab.NewClient(token)
}
