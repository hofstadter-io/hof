package git

import (
	"github.com/go-git/go-billy/v5"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage"
)

type GitRepo struct {
	Store storage.Storer
	FS    billy.Filesystem

	Repo *gogit.Repository

	Remote *gogit.Remote

	FetchOptions *gogit.FetchOptions
	ListOptions  *gogit.ListOptions
}
