package dag

import (
	"fmt"
	"regexp"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hashicorp/go-envparse"
	"github.com/hofstadter-io/hof/lib/env"
)

func (d *Dag) stepEnvVarHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	// no type for this one, it's just a map
	d.mx.RLock()
	var envs map[string]string
	err := step.Decode(&envs)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepEnv: %w", err)
	}

	for k, v := range envs {
		if k != "$kind" {
			c = c.WithEnvVariable(k, v, dagger.ContainerWithEnvVariableOpts{
				Expand: true,
			})
		}
	}
	return c, nil
}

type stepEnvVarsConfig struct {
	Kind string    `json:"$kind"`
	File cue.Value `json:"file"`
}

// needed for checking below, loop copied from module source
var envpair envparse.Pair

func (d *Dag) stepEnvFileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {

	// DEV HACK
	// return c, nil

	d.mx.RLock()
	var cfg stepEnvVarsConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepEnvfile: %w", err)
	}
	k := cfg.File.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return c, fmt.Errorf("missing $kind in stepEnvfile source: %v, got %v", step, k)
	}

	var file *dagger.File
	ks, _ := k.String()
	switch ks {
	case "#file":
		file, _, err = d.hashFile(cfg.File)

	case "#hostFile":
		file, _, err = d.HashHostFile(cfg.File)

	default:
		return c, fmt.Errorf("unsupported $kind in envfile.file: %v", step)
	}

	if err != nil {
		return nil, err
	}

	contents, err := file.Contents(d.ctx)
	if err != nil {
		return nil, err
	}

	r := strings.NewReader(contents)
	parser := envparse.New(r)
	for {
		kv, err := parser.Next()
		if err != nil {
			return nil, err
		}

		if kv == envpair {
			break
		}

		c = c.WithEnvVariable(kv.Key, kv.Val)
	}

	return c, nil
}

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

func (idx *hashSecretIndex) Key() string {
	if idx.cfg == nil {
		return "#secret.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#exportFile.%s", mk)
	}
	return fmt.Sprintf("#secret.%s", idx.cfg.Name)
}

func (d *Dag) hashSecret(step cue.Value) (*dagger.Secret, error) {
	d.mx.RLock()
	var cfg hashSecretConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding hashSecret: %w", err)
	}

	// index for query and create if not found
	idx := &hashSecretIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
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
			file, _, err := d.hashFile(cfg.Source)
			if err != nil {
				return nil, err
			}
			text, err := file.Contents(d.ctx)
			if err != nil {
				return nil, err
			}
			idx.shh = d.dag.SetSecret(cfg.Name, text)

		case "#hostFile":
			file, _, err := d.HashHostFile(cfg.Source)
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
	d.cat.Store(idx, idx)

	return idx.shh, nil
}

func (d *Dag) stepSecretVarHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c, nil
}

func (d *Dag) stepSecretVarsHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c, nil
}
