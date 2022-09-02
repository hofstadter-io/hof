package repos

import (
	"github.com/go-git/go-billy/v5"
)

type Repo struct {
	origUrl string
	origRef string

	Repo   string
	Subdir string
	Ver    string
	
	modDir string

	FS billy.Filesystem
}

func NewRepo(repo, ref string) *Repo {
	R := &Repo{
		origUrl: repo,
		origRef: ref,
	}

	return R
}

func (R *Repo) Resolve() error {

	return nil
}

func (R *Repo) Load() error {

	return nil
}
