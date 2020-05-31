package git

import (
	"github.com/go-git/go-git/v5/plumbing"
)

func (R *GitRepo) RemoteRefs() ([]*plumbing.Reference, error) {
	return R.Remote.List(R.ListOptions)
}

func (R *GitRepo) ClonedRefs() ([]*plumbing.Reference, error) {
	riter, err := R.Repo.References()
	if err != nil {
		return nil, err
	}

	var refs []*plumbing.Reference
	err = riter.ForEach(func(ref *plumbing.Reference) error {
		refs = append(refs, ref)
		return nil
	})

	return refs, err
}
