package dag

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hashicorp/go-envparse"
	"github.com/hofstadter-io/hof/lib/env"
)

func (d *Dag) stepEnvVarsHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	// no type for this one, it's just a map
	d.mx.RLock()
	var envs map[string]string
	err := step.Decode(&envs)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepEnv: %w", err)
	}

	keys := make([]string, 0, len(envs))
	for k := range envs {
		if k != "$kind" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := envs[k]
		c = c.WithEnvVariable(k, v, dagger.ContainerWithEnvVariableOpts{
			Expand: true,
		})
	}
	return c, nil
}

func (d *Dag) stepEnvAllHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	envs := os.Environ()
	sort.Strings(envs)
	for _, env := range envs {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			c = c.WithEnvVariable(parts[0], parts[1])
		}
	}
	return c, nil
}

type stepEnvFileConfig struct {
	Kind string    `json:"$kind"`
	File cue.Value `json:"file"`
}

// needed for checking below, loop copied from module source
var envpair envparse.Pair

func (d *Dag) stepEnvFileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {

	d.mx.RLock()
	var cfg stepEnvFileConfig
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
		file, _, err = d.hashFile(cfg.File, false)

	case "#hostFile":
		file, _, err = d.HashHostFile(cfg.File, false)

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
		upper := strings.ToUpper(sv) == sv
		if upper {
			// assume ENV var
			val := os.Getenv(sv)
			idx.shh = d.dag.SetSecret(cfg.Name, val)
		} else if matched {
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
			file, _, err := d.hashFile(cfg.Source, false)
			if err != nil {
				return nil, err
			}
			text, err := file.Contents(d.ctx)
			if err != nil {
				return nil, err
			}
			idx.shh = d.dag.SetSecret(cfg.Name, text)

		case "#hostFile":
			file, _, err := d.HashHostFile(cfg.Source, false)
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
	if idx.shh == nil {
		idx.shh = d.dag.Secret(cfg.Name, dagger.SecretOpts{})
	}
	// memoize
	d.cat.Store(idx, idx)

	return idx.shh, nil
}

func (d *Dag) stepSecretVarsHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	// no type for this one, it's just a map
	d.mx.RLock()
	var envs map[string]cue.Value
	err := step.Decode(&envs)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepSecretVar: %w", err)
	}

	keys := make([]string, 0, len(envs))
	for k := range envs {
		if k != "$kind" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := envs[k]
		s, err := d.hashSecret(v)
		if err != nil {
			return nil, fmt.Errorf("while decoding stepSecretVar.%s: %w", k, err)
		}
		c = c.WithSecretVariable(k, s)
	}
	return c, nil
}

func (d *Dag) stepSecretFileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepEnvFileConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepSecretFile: %w", err)
	}
	k := cfg.File.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return c, fmt.Errorf("missing $kind in stepSecretFile source: %v, got %v", step, k)
	}

	var contents string
	var file *dagger.File
	ks, _ := k.String()
	switch ks {
	case "#file":
		file, _, err = d.hashFile(cfg.File, false)
		if err != nil {
			break
		}
		contents, err = file.Contents(d.ctx)
		if err != nil {
			break
		}

	case "#hostFile":
		file, _, err = d.HashHostFile(cfg.File, false)
		if err != nil {
			break
		}
		contents, err = file.Contents(d.ctx)
		if err != nil {
			break
		}

	case "#secret":
		var s *dagger.Secret
		s, err = d.hashSecret(cfg.File)
		if err != nil {
			break
		}
		contents, err = s.Plaintext(d.ctx)
		if err != nil {
			break
		}

	default:
		return c, fmt.Errorf("unsupported $kind in stepSecretFile.file: %v", step)
	}

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

		s := d.dag.SetSecret(kv.Key, kv.Val)
		c = c.WithSecretVariable(kv.Key, s)
	}

	return c, nil
}
