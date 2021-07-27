package gitlab

import (
	"os"

	"github.com/xanzy/go-gitlab"
)

func NewClient() (client *gitlab.Client, err error) {
	return gitlab.NewClient(os.Getenv("GITLAB_TOKEN"))
}
