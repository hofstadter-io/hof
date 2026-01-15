package dag

import (
	"fmt"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

// HMMM, we don't seem to be using the index / memoization in this file
//       we do need to be careful about unnamed things, skip memo-optimization for those

type hashFileConfig struct {
	Kind       string    `json:"$kind"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	TrimPrefix string    `json:"trimPrefix"`
	Source     cue.Value `json:"source"`
}

type hashFileIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashFileConfig
	file *dagger.File
}

func (idx *hashFileIndex) Key() string {
	if idx.cfg == nil {
		return "#file.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#file.%s", mk)
	}
	return fmt.Sprintf("#file.%s", idx.cfg.Path)
}

func (d *Dag) hashFile(step cue.Value, noCache bool) (*dagger.File, string, error) {
	var err error
	step, err = d.ResolveShouldi(step)
	if err != nil {
		return nil, "", err
	}
	if !step.Exists() {
		return nil, "", fmt.Errorf("hashFile: resolved to empty value")
	}

	d.mx.RLock()
	var cfg hashFileConfig
	err = step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, "", fmt.Errorf("while decoding hashFile: %w", err)
	}

	// index for query and create if not found
	idx := &hashFileIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashFileIndex)
		return ix.file, ix.cfg.Path, nil
	}

	var (
		f   *dagger.File
		dir *dagger.Directory
		ctr *dagger.Container
	)

	sk := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !sk.Exists() {
		return nil, "", fmt.Errorf("missing $kind in struct file source: %v", step)
	}
	sks, _ := sk.String()
	switch sks {
	case "#gitRepo":
		repo, rcfg, err := d.hashGitRepo(cfg.Source, noCache)
		if err != nil {
			return nil, "", err
		}
		if rcfg != nil && rcfg.Ref != "" {
			dir = repo.Ref(rcfg.Ref).Tree()
		} else {
			dir = repo.Head().Tree()
		}
		f = dir.File(cfg.Path)

	case "#dir":
		dir, _, err = d.hashDir(cfg.Source, noCache)
		if err != nil {
			return nil, "", err
		}

		f = dir.File(cfg.Path)

	case "#hostDir":
		dir, _, err = d.HashHostDir(cfg.Source, noCache)
		if err != nil {
			return nil, "", err
		}
		f = dir.File(cfg.Path)

	case "#file":
		f, _, err = d.hashFile(cfg.Source, noCache)

	case "#hostFile":
		f, _, err = d.HashHostFile(cfg.Source, noCache)

	case "#rootfs":
		dir, err = d.hashRootFS(cfg.Source, noCache)
		if err == nil && dir != nil {
			f = dir.File(cfg.Path)
		}

	case "#cuefigSBOM":
		f, _, err = d.HashCuefigSBOM(cfg.Source, noCache)

	// TODO, make similar FileLike and ImageLike handlers so we don't repeat this everywhere
	case "#container":
		ctr, err = d.HashContainer(cfg.Source, noCache)
		if err != nil {
			return nil, "", err
		}
		f = ctr.File(cfg.Path)

	case "#hostImage":
		ctr, err = d.HashHostImage(cfg.Source, noCache)
		if err != nil {
			return nil, "", err
		}
		f = ctr.File(cfg.Path)

	case "#dockerBuild":
		ctr, err = d.HashDockerBuild(cfg.Source, noCache)
		if err != nil {
			return nil, "", err
		}
		f = ctr.File(cfg.Path)

	default:
		return nil, "", fmt.Errorf("hashFile.source: unsupported $kind %q in %v", sks, step)
	}

	if f == nil {
		return nil, "", fmt.Errorf("error hashFile.file result is nil in: %v", step)
	}

	// memoize
	idx.file = f
	d.cat.Store(idx, idx)

	path := idx.cfg.Path
	if cfg.TrimPrefix != "" {
		path = strings.TrimPrefix(path, cfg.TrimPrefix)
	}

	return idx.file, path, nil
}

type hashDirConfig struct {
	Kind    string      `json:"$kind"`
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	Sources []cue.Value `json:"sources"`

	TrimPrefix string    `json:"trimPrefix"`
	Patch      cue.Value `json:"patch"`
	PatchFile  cue.Value `json:"patchFile"`

	Include   []string `json:"include"`
	Exclude   []string `json:"exclude"`
	Gitignore bool     `json:"gitignore"`
}

type hashDirIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashDirConfig
	dir  *dagger.Directory
}

func (idx *hashDirIndex) Key() string {
	if idx.cfg == nil {
		return "#dir.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#dir.%s", mk)
	}
	return fmt.Sprintf("#dir.%s", idx.cfg.Path)
}

func (d *Dag) hashDir(step cue.Value, noCache bool) (*dagger.Directory, string, error) {
	var err error
	step, err = d.ResolveShouldi(step)
	if err != nil {
		return nil, "", err
	}
	if !step.Exists() {
		return nil, "", fmt.Errorf("hashDir: resolved to empty value")
	}

	d.mx.RLock()
	var cfg hashDirConfig
	err = step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, "", err
	}

	// index for query and create if not found
	idx := &hashDirIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashDirIndex)
		return ix.dir, ix.cfg.Path, nil
	}

	if len(cfg.Sources) == 0 {
		return nil, "", fmt.Errorf("empty sources decoding hashDir(%s): %w", cfg.Name, err)
	}

	// TODO, we need to do something similar for #Dir as we do here (bundle, multi-source)
	bundle := d.dag.Directory()
	for i, src := range cfg.Sources {
		d.mx.RLock()
		var k kinder
		err := src.Decode(&k)
		d.mx.RUnlock()
		if err != nil {
			return nil, "", fmt.Errorf("while decoding hashDir(%s).source.%d.$kind: %w", cfg.Name, i, err)
		}
		var (
			file *dagger.File
			dir  *dagger.Directory
			path string
		)
		switch k.Kind {
		case "#file":
			file, path, err = d.hashFile(src, noCache)
		case "#hostFile":
			_file, _, _err := d.HashHostFile(src, noCache)
			if _err == nil {
				file, err = _file, _err
			} else {
				err = _err
			}

		case "#dir":
			dir, path, err = d.hashDir(src, noCache)
		case "#hostDir":
			_dir, _, _err := d.HashHostDir(src, noCache)
			if _err == nil {
				dir, err = _dir, _err
			} else {
				err = _err
			}
		case "#gitRepo":
			repo, rcfg, rerr := d.hashGitRepo(src, noCache)
			if rerr == nil {
				if rcfg != nil && rcfg.Ref != "" {
					dir = repo.Ref(rcfg.Ref).Tree()
				} else {
					dir = repo.Head().Tree()
				}
			} else {
				err = rerr
			}

		case "#rootfs":
			dir, err = d.hashRootFS(src, noCache)

		case "#container":
			_ctr, _err := d.HashContainer(src, noCache)
			if _err == nil {
				dir = _ctr.Directory("/")
				path = "/"
			} else {
				err = _err
			}

		case "#hostImage":
			_ctr, _err := d.HashHostImage(src, noCache)
			if _err == nil {
				dir = _ctr.Directory("/")
			} else {
				err = _err
			}

		case "#dockerBuild":
			_ctr, _err := d.HashDockerBuild(src, noCache)
			if _err == nil {
				dir = _ctr.Directory("/")
			} else {
				err = _err
			}

		default:
			return nil, "", fmt.Errorf("unsupported kind %q in hashDir.source.%d.$kind", k.Kind, i)

		}
		if err != nil {
			return nil, "", fmt.Errorf("while decoding hashDir.source.%d.$kind: %w", i, err)
		}

		if file != nil {
			bundle = bundle.WithFile(path, file)
		}
		if dir != nil {
			bundle = bundle.WithDirectory(path, dir)
		}

	}

	// our bundle is assembled, craft the final dir
	// (1) filters
	final := d.dag.Directory().WithDirectory("/", bundle, dagger.DirectoryWithDirectoryOpts{
		Include:   cfg.Include,
		Exclude:   cfg.Exclude,
		Gitignore: cfg.Gitignore,
	})

	// (2) subpath selections
	if cfg.TrimPrefix != "" {
		final = final.Directory(cfg.TrimPrefix)
	}

	if cfg.Patch.Exists() {
		switch ik := cfg.Patch.IncompleteKind(); ik {
		case cue.StringKind:
			s, _ := cfg.Patch.String()
			final = final.WithPatch(s)
		case cue.StructKind:
			chg, err := d.HashChanges(cfg.Patch, noCache)
			if err != nil {
				return nil, "", err
			}
			f := chg.AsPatch()
			final = final.WithPatchFile(f)
		}
	} else if cfg.PatchFile.Exists() {
		f, err := d.HashPatchFile(cfg.PatchFile, noCache)
		if err != nil {
			return nil, "", err
		}
		final = final.WithPatchFile(f)
	}

	// memoize
	idx.dir = final
	d.cat.Store(idx, idx)

	return idx.dir, idx.cfg.Path, nil
}

type stepFileConfig struct {
	Kind string `json:"$kind"`
	// args
	Path    string    `json:"path"`
	Content cue.Value `json:"content"`
	// opts
	Permissions int    `json:"permissions"`
	Owner       string `json:"owner"`
	Expand      bool   `json:"expand"`
}

func (d *Dag) stepFileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepFileConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return c, err
	}

	if !cfg.Content.Exists() {
		return nil, fmt.Errorf("missing 'content' field in env.File: %v", step)
	}

	var f *dagger.File
	var dir *dagger.Directory
	switch ik := cfg.Content.IncompleteKind(); ik {
	case cue.StringKind:
		_, name := filepath.Split(cfg.Path)
		s, _ := cfg.Content.String()
		// fmt.Println("string file content:", cfg.Path)
		// fmt.Println(s)
		f = d.dag.File(name, s)
	case cue.StructKind:
		// look for kind
		k := cfg.Content.LookupPath(cue.ParsePath("$kind"))
		if !k.Exists() {
			return c, fmt.Errorf("missing $kind in struct file source: %v", step)
		}
		ks, _ := k.String()
		switch ks {
		case "#file":
			f, _, err = d.hashFile(cfg.Content, false) // steps should not use noCache from DAG?
		case "#hostFile":
			_file, _, _err := d.HashHostFile(cfg.Content, false)
			f, err = _file, _err
		case "#cuefigSBOM":
			f, _, err = d.HashCuefigSBOM(cfg.Content, false)

		case "#dir":
			dir, _, err = d.hashDir(cfg.Content, false)
			if err == nil && dir != nil {
				f = dir.File(cfg.Path)
			}
		case "#hostDir":
			dir, _, err = d.HashHostDir(cfg.Content, false)
			if err == nil && dir != nil {
				f = dir.File(cfg.Path)
			}

		case "#rootfs":
			dir, err = d.hashRootFS(cfg.Content, false)
			if err == nil && dir != nil {
				f = dir.File(cfg.Path)
			}

		default:
			return c, fmt.Errorf("unsupported $kind %q in struct file source: %v", ks, step)
		}

	default:
		return c, fmt.Errorf("unhandle incomplete cue kind %q in in struct file source: %v", ik, step)
	}

	if err != nil {
		return nil, fmt.Errorf("while trying to get file for content in %q: %w", cfg.Path, err)
	}

	if f == nil {
		return nil, fmt.Errorf("ERROR! should not get here, nil file from cue: %v", step)
	}
	c = c.WithFile(cfg.Path, f)

	return c, nil
}

type stepDirConfig struct {
	Kind string `json:"$kind"`
	// args
	Path   string    `json:"path"`
	Source cue.Value `json:"source"`
	// opts
	Include    []string  `json:"include"`
	Exclude    []string  `json:"exclude"`
	TrimPrefix string    `json:"trimPrefix"`
	Gitignore  bool      `json:"gitignore"`
	Owner      string    `json:"owner"`
	Expand     bool      `json:"expand"`
	Patch      cue.Value `json:"patch"`
	PatchFile  cue.Value `json:"patchFile"`
}

func (d *Dag) stepDirHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepDirConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepDir: %w", err)
	}

	var dir *dagger.Directory
	switch ik := cfg.Source.IncompleteKind(); ik {
	case cue.StructKind:
		// look for kind
		k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
		if !k.Exists() {
			return c, fmt.Errorf("missing $kind in stepDir source: %v", step)
		}
		ks, _ := k.String()
		switch ks {
		case "#dir":
			dir, _, err = d.hashDir(cfg.Source, false)
			if err != nil {
				return nil, err
			}
			// dir = dir.Directory(cfg.Path)

		case "#hostDir":
			dir, _, err = d.HashHostDir(cfg.Source, false)
			if err != nil {
				return nil, err
			}
			// dir = dir.Directory(cfg.Path)

		case "#gitRepo":
			repo, rcfg, err := d.hashGitRepo(cfg.Source, false)
			if err != nil {
				return nil, err
			}
			if rcfg != nil && rcfg.Ref != "" {
				dir = repo.Ref(rcfg.Ref).Tree()
			} else {
				dir = repo.Head().Tree()
			}

		case "#container":
			ctr, err := d.HashContainer(cfg.Source, false)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory("/")
			// dir = ctr.Directory(cfg.Path)

		case "#hostImage":
			ctr, err := d.HashHostImage(cfg.Source, false)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory("/")
			// dir = ctr.Directory(cfg.Path)

		case "#dockerBuild":
			ctr, err := d.HashDockerBuild(cfg.Source, false)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory("/")
			// dir = ctr.Directory(cfg.Path)

		case "#rootfs":
			dir, err = d.hashRootFS(cfg.Source, false)
			if err != nil {
				return nil, err
			}

		default:
			return c, fmt.Errorf("unsupported $kind in stepDir source: %v", step)
		}

	default:
		return c, fmt.Errorf("unsupported stepDir value type: %v", step)
	}

	if cfg.TrimPrefix != "" {
		dir = dir.Directory(cfg.TrimPrefix)
	}

	if cfg.Patch.Exists() {
		switch ik := cfg.Patch.IncompleteKind(); ik {
		case cue.StringKind:
			s, _ := cfg.Patch.String()
			dir = dir.WithPatch(s)
		case cue.StructKind:
			chg, err := d.HashChanges(cfg.Patch, false)
			if err != nil {
				return nil, err
			}
			f := chg.AsPatch()
			dir = dir.WithPatchFile(f)
		}
	} else if cfg.PatchFile.Exists() {
		f, err := d.HashPatchFile(cfg.PatchFile, false)
		if err != nil {
			return nil, err
		}
		dir = dir.WithPatchFile(f)
	}

	c = c.WithDirectory(cfg.Path, dir, dagger.ContainerWithDirectoryOpts{
		Include:   cfg.Include,
		Exclude:   cfg.Exclude,
		Gitignore: cfg.Gitignore,
		Owner:     cfg.Owner,
		Expand:    cfg.Expand,
	})
	return c, nil
}

type hashRootFSConfig struct {
	Kind   string    `json:"$kind"`
	Source cue.Value `json:"source"`
}

type hashRootFSIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashRootFSConfig
	dir  *dagger.Directory
}

func (idx *hashRootFSIndex) Key() string {
	if idx.cfg == nil {
		return "#rootfs.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#rootfs.%s", mk)
	}
	return fmt.Sprintf("#rootfs.%v", idx.val)
}

func (d *Dag) hashRootFS(val cue.Value, noCache bool) (*dagger.Directory, error) {
	var err error
	val, err = d.ResolveShouldi(val)
	if err != nil {
		return nil, err
	}
	if !val.Exists() {
		return nil, fmt.Errorf("hashRootFS: resolved to empty value")
	}

	d.mx.RLock()
	var cfg hashRootFSConfig
	err = val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding hashRootFS: %w", err)
	}

	// index for query and create if not found
	idx := &hashRootFSIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashRootFSIndex)
		return ix.dir, nil
	}

	var ctr *dagger.Container
	sk := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !sk.Exists() {
		return nil, fmt.Errorf("missing $kind in #rootfs source: %v", val)
	}
	sks, _ := sk.String()
	switch sks {
	case "#container":
		ctr, err = d.HashContainer(cfg.Source, noCache)
	case "#hostImage":
		ctr, err = d.HashHostImage(cfg.Source, noCache)
	case "#dockerBuild":
		ctr, err = d.HashDockerBuild(cfg.Source, noCache)
	default:
		return nil, fmt.Errorf("hashRootFS.source: unsupported $kind %q in %v", sks, val)
	}
	if err != nil {
		return nil, err
	}

	idx.dir = ctr.Rootfs()
	d.cat.Store(idx, idx)

	return idx.dir, nil
}

type stepRootFSConfig struct {
	Kind   string    `json:"$kind"`
	Source cue.Value `json:"source"`
}

func (d *Dag) stepRootFSHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepRootFSConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepRootFS: %w", err)
	}

	var dir *dagger.Directory
	sk := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !sk.Exists() {
		return nil, fmt.Errorf("missing $kind in RootFS source: %v", step)
	}
	sks, _ := sk.String()
	switch sks {
	case "#dir":
		dir, _, err = d.hashDir(cfg.Source, false)
	case "#hostDir":
		dir, _, err = d.HashHostDir(cfg.Source, false)
	case "#gitRepo":
		repo, rcfg, _err := d.hashGitRepo(cfg.Source, false)
		if _err == nil {
			if rcfg != nil && rcfg.Ref != "" {
				dir = repo.Ref(rcfg.Ref).Tree()
			} else {
				dir = repo.Head().Tree()
			}
		} else {
			err = _err
		}
	case "#rootfs":
		dir, err = d.hashRootFS(cfg.Source, false)
	case "#container":
		var ctr *dagger.Container
		ctr, err = d.HashContainer(cfg.Source, false)
		if err == nil {
			dir = ctr.Rootfs()
		}
	case "#hostImage":
		var ctr *dagger.Container
		ctr, err = d.HashHostImage(cfg.Source, false)
		if err == nil {
			dir = ctr.Rootfs()
		}
	case "#dockerBuild":
		var ctr *dagger.Container
		ctr, err = d.HashDockerBuild(cfg.Source, false)
		if err == nil {
			dir = ctr.Rootfs()
		}
	default:
		return nil, fmt.Errorf("RootFS.source: unsupported $kind %q in %v", sks, step)
	}

	if err != nil {
		return nil, err
	}

	return c.WithRootfs(dir), nil
}
