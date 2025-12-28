package dag

import (
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type gitRepoConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`
	Url  string `json:"url"`

	// opts
	KeepGitDir              bool      `json:"keepGitDir"`
	SSHKnownHosts           string    `json:"sshKnownHosts"`
	SSHAuthSocket           cue.Value `json:"sshAuthSocket"` // host socket
	HTTPAuthUsername        string    `json:"httpAuthUsername"`
	HTTPAuthToken           cue.Value `json:"httpAuthToken"`           // secret
	HTTPAuthHeader          cue.Value `json:"httpAuthHeader"`          // secret
	ExperimentalServiceHost cue.Value `json:"experimentalServiceHost"` // service
}

type gitRepoIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *gitRepoConfig
	repo *dagger.GitRepository
}

func (idx *gitRepoIndex) Key() string {
	if idx.cfg == nil {
		return "#gitRepo.nil"
	}
	return fmt.Sprintf("#gitRepo.%s", idx.cfg.Name)
}

func (d *Dag) hashGitRepo(step cue.Value) (*dagger.GitRepository, error) {
	var cfg gitRepoConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashHostFile: %w", err)
	}

	// index for query and create if not found
	idx := &gitRepoIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*gitRepoIndex)
		return ix.repo, nil
	}

	// load for realz
	idx.repo = d.dag.Git(cfg.Url, dagger.GitOpts{
		KeepGitDir:       cfg.KeepGitDir,
		SSHKnownHosts:    cfg.SSHKnownHosts,
		HTTPAuthUsername: cfg.HTTPAuthUsername,

		// the rest need to be decoded if they exist
		// SSHAuthSocket:  cfg.SSHAuthSocket,
		// HTTPAuthHeader: cfg.HTTPAuthHeader,
		// HTTPAuthToken: cfg.HTTPAuthToken,
		// ExperimentalServiceHost: cfg.ExperimentalServiceHost,
	})

	// memoize
	d.cat[idx] = idx

	return idx.repo, nil
}
