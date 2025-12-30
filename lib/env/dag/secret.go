package dag

import (
	"fmt"
	"regexp"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashSecretConfig struct {
	Kind   string    `json:"$kind"`
	Name   string    `json:"name"`
	Source cue.Value `json:"source"`
}

type hashSecretIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashSecretConfig
	shh  *dagger.Secret
}

func (h *hashSecretIndex) Key() string {
	if h.cfg == nil {
		return "#secret.nil"
	}
	return fmt.Sprintf("#secret.%s", h.cfg.Name)
}

func (d *Dag) hashSecret(step cue.Value) (*dagger.Secret, error) {
	var cfg hashSecretConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashSecret: %w", err)
	}

	// index for query and create if not found
	idx := &hashSecretIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hashSecretIndex)
		return ix.shh, nil
	}

	// load for realz
	// what kind of secret?
	switch ik := cfg.Source.IncompleteKind(); ik {
	case cue.StringKind:
		sv, _ := cfg.Source.String()
		re := regexp.MustCompile(`^[a-z]+://.*`)
		matched := re.MatchString(sv)
		if matched {
			// uri style secret for dagger
			idx.shh = d.dag.Secret(sv, dagger.SecretOpts{
				CacheKey: cfg.Name,
			})
		} else {
			// plaintext?!
			idx.shh = d.dag.SetSecret(cfg.Name, sv)
		}

	case cue.StructKind:
		kv := cfg.Source.LookupPath(cue.ParsePath("$kind"))
		if !kv.Exists() {
			return nil, fmt.Errorf("missing $kind in secret.source struct")
		}
		k, _ := kv.String()

		switch k {
		case "#file":
			file, err := d.hashFile(cfg.Source)
			if err != nil {
				return nil, err
			}
			text, err := file.Contents(d.ctx)
			if err != nil {
				return nil, err
			}
			idx.shh = d.dag.SetSecret(cfg.Name, text)

		case "#hostFile":
			file, err := d.hashHostFile(cfg.Source)
			if err != nil {
				return nil, err
			}
			text, err := file.Contents(d.ctx)
			if err != nil {
				return nil, err
			}
			idx.shh = d.dag.SetSecret(cfg.Name, text)
		default:
			return nil, fmt.Errorf("unsupported secret.source kind: %v", k)
		}

	default:
		return nil, fmt.Errorf("unsupported secret.source type: %v", ik)
	}
	idx.shh = d.dag.Secret(cfg.Name, dagger.SecretOpts{})
	// memoize
	d.cat[idx] = idx

	return idx.shh, nil
}

func (d *Dag) stepSecretHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c, nil
}
