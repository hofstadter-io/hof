package remote

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/lib/repos/git"
	"github.com/hofstadter-io/hof/lib/repos/oci"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)

// Parse parses a module name and returns
// the appropriate remote for it.
func Parse(mod string) (*Remote, error) {
	// TODO: Should pass a context in.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	k, err := NewKnowns()
	if err != nil {
		return nil, fmt.Errorf("new knowns: %w", err)
	}

	r := Remote{
		mod:    mod,
		knowns: k,
	}

	r.Host, r.Owner, r.Name = utils.ParseModURL(mod)

	// TODO: Store knowns

	isGit, err := git.IsGit(ctx, r.Host, r.Owner, r.Name)
	switch {
	case err != nil:
		return nil, fmt.Errorf("is git: %w", err)
	case isGit:
		r.kind = KindGit
		return &r, nil
	case oci.IsOCI(mod):
		r.kind = KindOCI
		return &r, nil
	}

	return nil, errors.New("remote not known")
}

type (
	LocalDir string
	Version  string
)

type Remote struct {
	Host  string
	Owner string
	Name  string

	mod    string
	kind   Kind
	knowns *Knowns
}

func (r *Remote) Pull(ctx context.Context, dir, ver string) error {
	switch r.kind {
	case KindGit:
		if err := git.SyncSource(dir, r.Host, r.Owner, r.Name, ver); err != nil {
			return fmt.Errorf("git sync source: %w", err)
		}
	case KindOCI:
		if err := oci.Pull(r.mod, dir); err != nil {
			return fmt.Errorf("oci pull: %w", err)
		}
	}

	return errors.New("pull: invalid kind")
}

type Kind string

const (
	KindGit Kind = "git"
	KindOCI Kind = "oci"
)
