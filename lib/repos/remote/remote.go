package remote

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hofstadter-io/hof/lib/repos/git"
	"github.com/hofstadter-io/hof/lib/repos/oci"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)


// Parse parses a module name and returns
// the appropriate remote for it.
func Parse(mod string) (*Remote, error) {


	// TODO: Is is worth having a complex type here to handle
	// this kind of check? Will Mirrors get used elsewhere?
	// If not, it would be better to simplify Mirrors to a
	// function.
	// defer m.Close()

	r := NewRemote(mod, MirrorsSingleton)
	_, err := r.Kind()
	if err != nil {
		return nil, err
	}

	return r, nil
}

type Remote struct {
	Host  string
	Owner string
	Name  string

	mod     string
	kind    Kind
	mirrors *Mirrors
}

func NewRemote(mod string, mir *Mirrors) *Remote {
	r := &Remote{
		mod:     mod,
		mirrors: mir,
	}

	r.Host, r.Owner, r.Name = utils.ParseModURL(mod)
	return r
}

func (r *Remote) Pull(ctx context.Context, dir, ver string) error {
	switch r.kind {
	case KindGit:
		// ensure we have up-to-date code in .cache/hof/src/<module>
		if err := git.SyncSource(dir, r.Host, r.Owner, r.Name); err != nil {
			return fmt.Errorf("git sync source: %w", err)
		}

	case KindOCI:
		// extract hash from version
		if strings.HasPrefix(ver, "v0.0.0-") {
			parts := strings.Split(ver, "-")
			ver = parts[len(parts)-1]
			// HACK
			ver = "sha256:" + ver
		}
		if err := oci.Pull(r.mod + "@" + ver, dir); err != nil {
			return fmt.Errorf("oci pull: %w", err)
		}

	default:
		return fmt.Errorf("usupported kind: %s", r.kind)
	}

	return nil
}

func (r *Remote) Publish(ctx context.Context, dir string, tag string) error {
	switch r.kind {
	case KindOCI:
		codeDir, err := oci.NewCode(dir)
		if err != nil {
			return fmt.Errorf("oci new code: %w", err)
		}

		img, err := oci.Build(dir, []oci.Dir{oci.NewDeps(), codeDir})
		if err != nil {
			return fmt.Errorf("oci build: %w", err)
		}

		if err := oci.Push(tag, img); err != nil {
			return fmt.Errorf("oci publish: %w", err)
		}

		return nil
	}

	return fmt.Errorf("unsupported kind: %s", r.kind)
}

func (r *Remote) Kind() (Kind, error) {
	if r.kind != "" {
		return  r.kind, nil
	}

	// TODO: Should pass a context in.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	isOCI, err := r.mirrors.Is(ctx, KindOCI, r.mod)
	if debug > 0 {
		// fmt.Println("isOCI?:", isOCI, err)
	}
	if isOCI && err != nil {
		return "", fmt.Errorf("mirror is oci: %w", err)
	}
	if isOCI {
		r.kind = KindOCI
		return "", nil
	}

	isGit, err := r.mirrors.Is(ctx, KindGit, r.mod)
	if debug > 0 {
		fmt.Println("isGit?:", isGit, err)
	}
	if isGit && err != nil {
		return "", fmt.Errorf("mirror is git: %w", err)
	}
	if isGit {
		r.kind = KindGit
		return "", nil
	}


	return "", fmt.Errorf("cannot determine registry kind for: %s", r.mod)
}

type Kind string

const (
	KindGit Kind = "git"
	KindOCI Kind = "oci"
)
