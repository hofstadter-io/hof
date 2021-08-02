package gitlab

import (
	"os"

	"github.com/xanzy/go-gitlab"
)

const TokenEnv = "GITLAB_TOKEN"

func NewClient() (client *gitlab.Client, err error) {
	return gitlab.NewClient(os.Getenv(TokenEnv))
}
