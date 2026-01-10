package dag

import (
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashGitRepoConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`
	Url  string `json:"url"`
	Ref  string `json:"ref"`

	// opts
	KeepGitDir              bool      `json:"keepGitDir"`
	SSHKnownHosts           string    `json:"sshKnownHosts"`
	SSHAuthSocket           cue.Value `json:"sshAuthSocket"` // host socket
	HTTPAuthUsername        string    `json:"httpAuthUsername"`
	HTTPAuthToken           cue.Value `json:"httpAuthToken"`           // secret
	HTTPAuthHeader          cue.Value `json:"httpAuthHeader"`          // secret
	ExperimentalServiceHost cue.Value `json:"experimentalServiceHost"` // service
}

type hashGitRepoIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashGitRepoConfig
	repo *dagger.GitRepository
}

// these index & Key are getting repetitive, as are the handlers
// we should be able to refactor these down with some generics & interface love
// we need to refactor down a bit anyway, use kinder in more places, more metadata inspection
// think about how the recursive CUE decoding and structs here will eventually
// be merged with and intermix with other subsystems like hof/flow

func (idx *hashGitRepoIndex) Key() string {
	if idx.cfg == nil {
		return "#gitRepo.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		// TODO, need a ref here, generally they may need parameteres
		return fmt.Sprintf("#gitRepo.%s", mk)
	}
	return fmt.Sprintf("#gitRepo.%s", idx.cfg.Name)
}

func (d *Dag) hashGitRepo(step cue.Value, noCache bool) (*dagger.GitRepository, *hashGitRepoConfig, error) {
	var err error
	step, err = d.ResolveShouldi(step)
	if err != nil {
		return nil, nil, fmt.Errorf("while resolving hashGitRepo: %w", err)
	}
	if !step.Exists() {
		return nil, nil, fmt.Errorf("hashGitRepo: resolved to empty value")
	}

	d.mx.RLock()
	var cfg hashGitRepoConfig
	err = step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostFile: %w", err)
	}

	// index for query and create if not found
	idx := &hashGitRepoIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashGitRepoIndex)
		return ix.repo, ix.cfg, nil
	}

	// fmt.Println("#GitRepo", step, cfg)

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
	d.cat.Store(idx, idx)

	return idx.repo, idx.cfg, nil
}
