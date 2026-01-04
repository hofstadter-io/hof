package dag

import (
	"fmt"
	"path/filepath"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

// HMMM, we don't seem to be using the index / memoization in this file
//       we do need to be careful about unnamed things, skip memo-optimization for those

type hashFileConfig struct {
	Kind   string    `json:"$kind"`
	Name   string    `json:"name"`
	Path   string    `json:"path"`
	Source cue.Value `json:"source"`
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

func (d *Dag) hashFile(step cue.Value) (*dagger.File, string, error) {
	var cfg hashFileConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, "", fmt.Errorf("while decoding hashFile: %w", err)
	}

	// index for query and create if not found
	idx := &hashFileIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
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
		repo, rcfg, err := d.hashGitRepo(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		dir = repo.Ref(rcfg.Ref).Tree()
		f = dir.File(cfg.Path)

	case "#dir":
		dir, _, err = d.hashDir(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		f = dir.File(cfg.Path)

	case "#hostDir":
		dir, _, err = d.HashHostDir(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		f = dir.File(cfg.Path)

	// TODO, make similar FileLike and ImageLike handlers so we don't repeat this everywhere
	case "#container":
		ctr, err = d.HashContainer(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		f = ctr.File(cfg.Path)

	case "#hostImage":
		ctr, err = d.HashHostImage(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		f = ctr.File(cfg.Path)

	case "#dockerBuild":
		ctr, err = d.HashDockerBuild(cfg.Source)
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
	d.cat[idx] = idx

	return idx.file, idx.cfg.Path, nil
}

type hashDirConfig struct {
	Kind    string      `json:"$kind"`
	Name    string      `json:"name"`
	Path    string      `json:"path"`
	Sources []cue.Value `json:"sources"`

	TrimPrefix string    `json:"trimPrefix"`
	Patch      string    `json:"patch"`
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

func (d *Dag) hashDir(step cue.Value) (*dagger.Directory, string, error) {
	var cfg hashDirConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, "", err
	}

	// index for query and create if not found
	idx := &hashDirIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
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
		var k kinder
		err := src.Decode(&k)
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
			file, path, err = d.hashFile(src)
		case "#hostFile":
			_file, _, _err := d.HashHostFile(src)
			if _err == nil {
				file, err = _file, _err
			} else {
				err = _err
			}

		case "#dir":
			dir, path, err = d.hashDir(src)
		case "#hostDir":
			_dir, _, _err := d.HashHostDir(src)
			if _err == nil {
				dir, err = _dir, _err
			} else {
				err = _err
			}
		case "#gitRepo":
			repo, rcfg, rerr := d.hashGitRepo(src)
			if rerr == nil {
				dir = repo.Ref(rcfg.Ref).Tree()
			} else {
				err = rerr
			}
		default:
			return nil, "", fmt.Errorf("unsupported kind %q in hashDir.source.%d.$kind: %w", k.Kind, i, err)

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

	if cfg.Patch != "" {
		final = final.WithPatch(cfg.Patch)
	} else if cfg.PatchFile.Exists() {
		f, _, err := d.hashFile(cfg.PatchFile)
		if err != nil {
			return nil, "", err
		}
		final = final.WithPatchFile(f)
	}

	// memoize
	idx.dir = final
	d.cat[idx] = idx

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
	var cfg stepFileConfig
	err := step.Decode(&cfg)
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
			f, _, err = d.hashFile(cfg.Content)
		case "#hostFile":
			_file, _, _err := d.HashHostFile(cfg.Content)
			f, err = _file, _err

		case "#dir":
			dir, _, err = d.hashDir(cfg.Content)
			if err == nil && dir != nil {
				f = dir.File(cfg.Path)
			}
		case "#hostDir":
			dir, _, err = d.HashHostDir(cfg.Content)
			f = dir.File(cfg.Path)
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
	Include   []string `json:"include"`
	Exclude   []string `json:"exclude"`
	Gitignore bool     `json:"gitignore"`
	Owner     string   `json:"owner"`
	Expand    bool     `json:"expand"`
}

func (d *Dag) stepDirHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepDirConfig
	err := step.Decode(&cfg)
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
			dir, _, err = d.hashDir(cfg.Source)
			if err != nil {
				return nil, err
			}
			// dir = dir.Directory(cfg.Path)

		case "#hostDir":
			dir, _, err = d.HashHostDir(cfg.Source)
			if err != nil {
				return nil, err
			}
			// dir = dir.Directory(cfg.Path)

		case "#gitRepo":
			repo, rcfg, err := d.hashGitRepo(cfg.Source)
			if err != nil {
				return nil, err
			}
			dir = repo.Ref(rcfg.Ref).Tree()

		case "#container":
			ctr, err := d.HashContainer(cfg.Source)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory("")
			// dir = ctr.Directory(cfg.Path)

		case "#hostImage":
			ctr, err := d.HashHostImage(cfg.Source)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory("")
			// dir = ctr.Directory(cfg.Path)

		case "#dockerBuild":
			ctr, err := d.HashDockerBuild(cfg.Source)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory("")
			// dir = ctr.Directory(cfg.Path)

		default:
			return c, fmt.Errorf("unsupported $kind in stepDir source: %v", step)
		}

	default:
		return c, fmt.Errorf("unsupported stepDir value type: %v", step)
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
