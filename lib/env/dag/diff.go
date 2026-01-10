package dag

import (
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/yagu"
)

type hashChangesConfig struct {
	Kind string    `json:"$kind"`
	Prev cue.Value `json:"prev"`
	Next cue.Value `json:"next"`
}

type hashChangesIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashChangesConfig
	chg  *dagger.Changeset
}

func (idx *hashChangesIndex) Key() string {
	if idx.cfg == nil {
		return "#changes.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#changes.%s", mk)
	}
	return "#changes"
}

func (d *Dag) HashChanges(val cue.Value, noCache bool) (*dagger.Changeset, error) {
	d.mx.RLock()
	var cfg hashChangesConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, err
	}

	idx := &hashChangesIndex{
		val: val,
		cfg: &cfg,
	}

	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashChangesIndex)
		return ix.chg, nil
	}

	// build
	var (
		prev *dagger.Directory
		next *dagger.Directory
	)

	// prev
	p, _, err := d.Dir(cfg.Prev, noCache)
	if err != nil {
		return nil, err
	}
	prev = p

	// next
	n, _, err := d.Dir(cfg.Next, noCache)
	if err != nil {
		return nil, err
	}
	next = n

	idx.chg = prev.Changes(next)
	d.cat.Store(idx, idx)

	return idx.chg, nil
}

type stepChangesConfig struct {
	Kind     string    `json:"$kind"`
	Change   cue.Value `json:"change"`
	Basepath string    `json:"basepath"`
}

func (d *Dag) stepChangesHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepChangesConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return c, err
	}

	chg, err := d.HashChanges(cfg.Change, false)
	if err != nil {
		return c, err
	}

	rfs := c.Rootfs().Directory(cfg.Basepath).WithChanges(chg)
	c = c.WithRootfs(rfs)

	return c, nil
}

type hashPatchFileConfig struct {
	Kind   string    `json:"$kind"`
	Source cue.Value `json:"source"`
	Name   string    `json:"name"`
}

type hashPatchFileIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashPatchFileConfig
	file *dagger.File
}

func (idx *hashPatchFileIndex) Key() string {
	if idx.cfg == nil {
		return "#patchFile.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#patchFile.%s", mk)
	}
	return "#patchFile"
}

func (d *Dag) HashPatchFile(val cue.Value, noCache bool) (*dagger.File, error) {
	d.mx.RLock()
	var cfg hashPatchFileConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, err
	}

	idx := &hashPatchFileIndex{
		val: val,
		cfg: &cfg,
	}

	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashPatchFileIndex)
		return ix.file, nil
	}

	name := cfg.Name
	if name == "" {
		name = "patch.diff"
	}

	var f *dagger.File
	switch ik := cfg.Source.IncompleteKind(); ik {
	case cue.StringKind:
		s, _ := cfg.Source.String()
		f = d.dag.File(name, s)
	case cue.StructKind:
		chg, cerr := d.HashChanges(cfg.Source, noCache)
		if cerr != nil {
			return nil, cerr
		}
		f = chg.AsPatch()
	default:
		return nil, fmt.Errorf("unsupported patch source kind: %v", ik)
	}

	idx.file = f
	d.cat.Store(idx, idx)

	return idx.file, nil
}

type stepPatchConfig struct {
	Kind     string    `json:"$kind"`
	Source   cue.Value `json:"source"`
	Basepath string    `json:"basepath"`
}

func (d *Dag) stepPatchHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepPatchConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return c, err
	}

	switch ik := cfg.Source.IncompleteKind(); ik {
	case cue.StringKind:
		s, _ := cfg.Source.String()
		rfs := c.Rootfs().WithPatch(s)
		c = c.WithRootfs(rfs)
		return c, nil
	case cue.StructKind:
		chg, err := d.HashChanges(cfg.Source, false)
		if err != nil {
			return c, err
		}
		rfs := c.Rootfs().Directory(cfg.Basepath).WithChanges(chg)
		c = c.WithRootfs(rfs)
		return c, nil
	default:
		return c, fmt.Errorf("unsupported patch source kind: %v", ik)
	}

}

type stepPatchFileConfig struct {
	Kind     string    `json:"$kind"`
	Source   cue.Value `json:"source"`
	Basepath string    `json:"basepath"`
}

func (d *Dag) stepPatchFileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepPatchFileConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return c, err
	}

	f, err := d.HashPatchFile(cfg.Source, false)
	if err != nil {
		return c, err
	}
	rfs := c.Rootfs().Directory(cfg.Basepath).WithPatchFile(f)
	c = c.WithRootfs(rfs)

	return c, nil
}

type hashShouldiConfig struct {
	Kind    string    `json:"$kind"`
	Changes cue.Value `json:"changes"`
	Include []string  `json:"include"`
	Exclude []string  `json:"exclude"`
	Force   bool      `json:"force"`
	Then    cue.Value `json:"then"`
	Else    cue.Value `json:"else"`
}

type hashShouldiIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashShouldiConfig
	// what does this return?
	res cue.Value
}

func (idx *hashShouldiIndex) Key() string {
	if idx.cfg == nil {
		return "#shouldi.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#shouldi.%s", mk)
	}
	return "#shouldi"
}

func (d *Dag) ResolveShouldi(val cue.Value) (cue.Value, error) {
	for {
		if !val.Exists() {
			return val, nil
		}

		// check for #shouldi
		var k kinder
		err := val.Decode(&k)
		if err != nil {
			// not a struct or matching shape, probably the value we want
			return val, nil
		}

		if k.Kind == "#shouldi" {
			next, err := d.HashShouldi(val)
			if err != nil {
				return val, err
			}
			val = next
			continue
		}

		break
	}
	return val, nil
}

func (d *Dag) HashShouldi(val cue.Value) (cue.Value, error) {
	d.mx.RLock()
	var cfg hashShouldiConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return val, err
	}

	idx := &hashShouldiIndex{
		val: val,
		cfg: &cfg,
	}

	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashShouldiIndex)
		return ix.res, nil
	}

	if cfg.Force {
		idx.res = cfg.Then
		d.cat.Store(idx, idx)
		return idx.res, nil
	}

	// Get the changes
	chg, err := d.HashChanges(cfg.Changes, false)
	if err != nil {
		return cue.Value{}, err
	}

	// Collect the changed paths
	addpaths, err := chg.AddedPaths(d.ctx)
	if err != nil {
		return cue.Value{}, err
	}
	modpaths, err := chg.ModifiedPaths(d.ctx)
	if err != nil {
		return cue.Value{}, err
	}
	delpaths, err := chg.RemovedPaths(d.ctx)
	if err != nil {
		return cue.Value{}, err
	}
	allpaths := make([]string, len(addpaths)+len(modpaths)+len(delpaths))
	allpaths = append(addpaths, modpaths...)
	allpaths = append(allpaths, delpaths...)

	// Calculate shouldi
	shouldi := false
	// Only proceed if there are any changes detected
	if len(allpaths) > 0 {
		// If no include or exclude patterns are specified, any change triggers shouldi
		if len(cfg.Include) == 0 && len(cfg.Exclude) == 0 {
			shouldi = true
		} else {
			// Check each changed path against patterns
			for _, p := range allpaths {
				inc, err := yagu.CheckShouldInclude(p, cfg.Include, cfg.Exclude)
				if err != nil {
					return cue.Value{}, err
				}
				if inc {
					shouldi = true
					break
				}
			}
		}
	}

	if shouldi {
		idx.res = cfg.Then
	} else {
		if cfg.Else.Exists() {
			idx.res = cfg.Else
		} else {
			idx.res = cue.Value{}
		}
	}

	d.cat.Store(idx, idx)
	return idx.res, nil
}
