package remote

import (
	"context"
	"errors"
	"fmt"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/hofstadter-io/hof/lib/repos/git"
	"github.com/hofstadter-io/hof/lib/repos/oci"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)

// Parse parses a module name and returns
// the appropriate remote for it.
func Parse(mod string) (Remote, error) {
	// TODO: Should pass a context in.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch {
	case isGit(ctx, mod):
		return newGitRemote(mod), nil
	case isOCI(ctx, mod):
		return ociRemote{mod: mod}, nil
	}

	return nil, errors.New("remote not known")
}

type (
	LocalDir string
	Version  string
)

type Remote interface {
	Parts() []string
	Pull(context.Context, LocalDir, Version) error
}

func newGitRemote(mod string) gitRemote {
	var r gitRemote
	r.host, r.owner, r.repo = utils.ParseModURL(mod)

	return r
}

type gitRemote struct {
	host  string
	owner string
	repo  string
}

func (r gitRemote) Pull(_ context.Context, d LocalDir, v Version) error {
	var (
		dir = string(d)
		ver = string(v)
	)

	if err := git.SyncSource(dir, r.host, r.owner, r.repo, ver); err != nil {
		return fmt.Errorf("git sync source: %w", err)
	}

	return nil
}

func (r gitRemote) Parts() []string {
	return []string{r.host, r.owner, r.repo}
}

func isGit(ctx context.Context, mod string) bool {
	// TODO: Cache these on disk once it is known what this is.
	rem := gogit.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://" + mod},
	})

	_, err := rem.ListContext(ctx, &gogit.ListOptions{})
	// TODO: This isn't ideal. This could be a failure
	// due to bad credentials and it would be better
	// to test for that and prompt the user.
	return err == nil
}

type ociRemote struct {
	mod string
}

// Parts implements Remote
func (r ociRemote) Parts() []string {
	return []string{r.mod}
}

// Pull implements Remote
func (r ociRemote) Pull(_ context.Context, d LocalDir, v Version) error {
	if err := oci.Pull(r.mod, string(d)); err != nil {
		return fmt.Errorf("oci pull: %w", err)
	}
	return nil
}

func isOCI(_ context.Context, mod string) bool {
	// TODO:
	//   * This is naive. Revisit it.
	//   * Cache these on disk once it is known what this is.
	_, err := crane.Head(mod)
	return err == nil
}
