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
	return fmt.Sprintf("#file.%s", idx.cfg.Path)
}

func (d *Dag) hashFile(step cue.Value) (*dagger.File, string, error) {
	var cfg hashFileConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, "", fmt.Errorf("while decoding hashFile: %w", err)
	}

	// ind

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
		dir := repo.Ref(rcfg.Ref).Tree()
		return dir.File(cfg.Path), cfg.Path, nil

	case "#dir":
		dir, _, err := d.hashDir(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		return dir.File(cfg.Path), cfg.Path, nil

	case "#hostDir":
		dir, _, err := d.hashHostDir(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		return dir.File(cfg.Path), cfg.Path, nil

	// TODO, make similar FileLike and ImageLike handlers so we don't repeat this everywhere
	case "#container":
		ctr, err := d.HashContainer(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		return ctr.File(cfg.Path), cfg.Path, nil

	case "#hostImage":
		ctr, err := d.HashHostImage(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		return ctr.File(cfg.Path), cfg.Path, nil

	case "#dockerBuild":
		ctr, err := d.HashDockerBuild(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		return ctr.File(cfg.Path), cfg.Path, nil

	default:
		return nil, "", fmt.Errorf("hashFile.source: unsupported $kind: %s", sks)
	}
}

type hashDirConfig struct {
	Kind      string    `json:"$kind"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Source    cue.Value `json:"source"`
	Patch     string    `json:"patch"`
	PatchFile cue.Value `json:"patchFile"`
	Include   []string  `json:"include"`
	Exclude   []string  `json:"exclude"`
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
	return fmt.Sprintf("#dir.%s", idx.cfg.Path)
}

func (d *Dag) hashDir(step cue.Value) (*dagger.Directory, string, error) {
	var cfg hashDirConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, "", err
	}

	sk := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !sk.Exists() {
		return nil, "", fmt.Errorf("missing $kind in struct file source: %v", step)
	}

	var dir *dagger.Directory
	sks, _ := sk.String()
	switch sks {
	case "#gitRepo":
		repo, rcfg, err := d.hashGitRepo(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		dir = repo.Ref(rcfg.Ref).Tree().Directory(cfg.Path)

	case "#dir":
		dir, _, err = d.hashDir(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		dir = dir.Directory(cfg.Path)

	case "#hostDir":
		dir, _, err = d.hashHostDir(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		dir = dir.Directory(cfg.Path)

	case "#container":
		ctr, err := d.HashContainer(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		dir = ctr.Directory(cfg.Path)

	case "#hostImage":
		ctr, err := d.HashHostImage(cfg.Source)
		if err != nil {
			return nil, "", err
		}
		dir = ctr.Directory(cfg.Path)

	default:
		return nil, "", fmt.Errorf("hashFile.source: unsupported $kind: %s", sks)
	}

	if cfg.Patch != "" {
		dir = dir.WithPatch(cfg.Patch)
	} else if cfg.PatchFile.Exists() {
		f, _, err := d.hashFile(cfg.PatchFile)
		if err != nil {
			return nil, "", err
		}
		dir = dir.WithPatchFile(f)
	}

	return dir, cfg.Path, nil
}

type hashExportFileConfig struct {
	Kind string    `json:"$kind"`
	Name string    `json:"name"`
	Path string    `json:"path"`
	File cue.Value `json:"file"`

	AllowParentDirPath bool `json:"allowParentDirPath"`
}

func (d *Dag) HashExportFile(step cue.Value) (*dagger.File, *hashExportFileConfig, error) {
	var cfg hashExportFileConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, nil, err
	}

	file, _, err := d.hashFile(cfg.File)
	if err != nil {
		return nil, nil, err
	}

	return file, &cfg, nil
}

type hashExportDirConfig struct {
	Kind string    `json:"$kind"`
	Name string    `json:"name"`
	Path string    `json:"path"`
	Dir  cue.Value `json:"dir"`
	Wipe bool      `json:"wipe"`
}

func (d *Dag) HashExportDir(step cue.Value) (*dagger.Directory, *hashExportDirConfig, error) {
	var cfg hashExportDirConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, nil, err
	}

	dir, _, err := d.hashDir(cfg.Dir)
	if err != nil {
		return nil, nil, err
	}

	return dir, &cfg, nil
}

type hashExportImageConfig struct {
	Kind  string    `json:"$kind"`
	Name  string    `json:"name"`
	Url   string    `json:"url"`
	Image cue.Value `json:"image"`
}

func (d *Dag) HashExportImage(step cue.Value) (*dagger.Container, *hashExportImageConfig, error) {
	var cfg hashExportImageConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, nil, err
	}

	c, err := d.HashContainer(cfg.Image)
	if err != nil {
		return nil, nil, err
	}

	return c, &cfg, nil
}

type hashExportImageFileConfig struct {
	Kind  string    `json:"$kind"`
	Name  string    `json:"name"`
	Path  string    `json:"path"`
	Image cue.Value `json:"image"`
}

func (d *Dag) HashExportImageFile(step cue.Value) (*dagger.Container, *hashExportImageFileConfig, error) {
	var cfg hashExportImageFileConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, nil, err
	}

	c, err := d.HashContainer(cfg.Image)
	if err != nil {
		return nil, nil, err
	}

	return c, &cfg, nil
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
			f, _, err = d.hashHostFile(cfg.Content)

		default:
			return c, fmt.Errorf("unsupported $kind in struct file source: %v", step)
		}
	}

	if f == nil {
		return nil, fmt.Errorf("ERROR! should not get here, nil file from cue", step)
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

		case "#hostDir":
			dir, _, err = d.hashHostDir(cfg.Source)
			if err != nil {
				return nil, err
			}

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
			dir = ctr.Directory(cfg.Path)

		case "#hostImage":
			ctr, err := d.HashHostImage(cfg.Source)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory(cfg.Path)

		// case "#dockerBuild":
		// 	ctr, err := d.HashDockerBuild(cfg.Source)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	dir = ctr.Directory(cfg.Path)

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
