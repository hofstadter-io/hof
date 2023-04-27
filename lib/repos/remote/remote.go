package remote

import (
	"context"
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

	m, err := NewMirrors()
	if err != nil {
		return nil, fmt.Errorf("new knowns: %w", err)
	}

	// TODO: Is is worth having a complex type here to handle
	// this kind of check? Will Mirrors get used elsewhere?
	// If not, it would be better to simplify Mirrors to a
	// function.
	defer m.Close()

	r := Remote{
		mod:     mod,
		mirrors: m,
	}

	r.Host, r.Owner, r.Name = utils.ParseModURL(mod)

	isOCI, err := m.Is(ctx, KindOCI, mod)
	if err != nil {
		return nil, fmt.Errorf("mirror is oci: %w", err)
	}
	if isOCI {
		r.kind = KindOCI
		return &r, nil
	}

	isGit, err := m.Is(ctx, KindGit, mod)
	if err != nil {
		return nil, fmt.Errorf("mirror is git: %w", err)
	}
	if isGit {
		r.kind = KindGit
		return &r, nil
	}

	return nil, fmt.Errorf("cannot parse %s", mod)
}

type Remote struct {
	Host  string
	Owner string
	Name  string

	mod     string
	kind    Kind
	mirrors *Mirrors
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

	return fmt.Errorf("usupported kind: %s", r.kind)
}

type Kind string

const (
	KindGit Kind = "git"
	KindOCI Kind = "oci"
)
