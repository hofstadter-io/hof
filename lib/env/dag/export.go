package dag

import (
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type exportFileConfig struct {
	Kind string    `json:"$kind"`
	Name string    `json:"name"`
	Path string    `json:"path"`
	File cue.Value `json:"file"`

	AllowParentDirPath bool `json:"allowParentDirPath"`
}

type exportFileIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *exportFileConfig
	file *dagger.File
}

func (idx *exportFileIndex) Key() string {
	if idx.cfg == nil {
		return "#exportFile.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#exportFile.%s", mk)
	}
	return fmt.Sprintf("#exportFile.%s", idx.cfg.Name)
}

func (d *Dag) HashExportFile(step cue.Value, noCache bool) (*dagger.File, *exportFileConfig, error) {
	d.mx.RLock()
	var cfg exportFileConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashExportFile: %w", err)
	}

	// index for query and create if not found
	idx := &exportFileIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*exportFileIndex)
		return ix.file, ix.cfg, nil
	}

	f, _, err := d.hashFile(cfg.File, noCache)
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashExportFile.file: %w", err)
	}

	// memoize
	idx.file = f
	d.cat.Store(idx, idx)

	return idx.file, idx.cfg, nil
}

type exportDirConfig struct {
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
	Owner     string   `json:"owner"`

	Wipe bool `json:"wipe"`
}

type exportDirIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *exportDirConfig
	dir  *dagger.Directory
}

func (idx *exportDirIndex) Key() string {
	if idx.cfg == nil {
		return "#exportDir.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#exportDir.%s", mk)
	}
	return fmt.Sprintf("#exportDir.%s", idx.cfg.Name)
}

func (d *Dag) HashExportDir(step cue.Value, noCache bool) (*dagger.Directory, *exportDirConfig, error) {
	d.mx.RLock()
	var cfg exportDirConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashExportDir: %w", err)
	}

	// index for query and create if not found
	idx := &exportDirIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*exportDirIndex)
		return ix.dir, ix.cfg, nil
	}

	if len(cfg.Sources) == 0 {
		return nil, nil, fmt.Errorf("empty sources decoding hashExportDir(%s): %w", cfg.Name, err)
	}

	// TODO, we need to do something similar for #Dir as we do here (bundle, multi-source)
	bundle := d.dag.Directory()
	for i, src := range cfg.Sources {
		d.mx.RLock()
		var k kinder
		err := src.Decode(&k)
		d.mx.RUnlock()
		if err != nil {
			return nil, nil, fmt.Errorf("while decoding hashExportDir(%s).source.%d.$kind: %w", cfg.Name, i, err)
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
			_file, _cfg, _err := d.HashHostFile(src, noCache)
			file, path, err = _file, _cfg.Path, _err

		case "#dir":
			dir, path, err = d.hashDir(src, noCache)
		case "#hostDir":
			_dir, _cfg, _err := d.HashHostDir(src, noCache)
			dir, path, err = _dir, _cfg.Path, _err
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

		case "#cuefigSBOM":
			_file, _path, _err := d.HashCuefigSBOM(src, noCache)
			file, path, err = _file, _path, _err

		default:
			return nil, nil, fmt.Errorf("unsupported kind %q in hashExportDir.source.%d.$kind: %w", k.Kind, i, err)

		}
		if err != nil {
			return nil, nil, fmt.Errorf("while decoding hashExportDir.source.%d.$kind: %w", i, err)
		}

		if file != nil {
			bundle = bundle.WithFile(path, file, dagger.DirectoryWithFileOpts{})
		}
		if dir != nil {
			bundle = bundle.WithDirectory(path, dir, dagger.DirectoryWithDirectoryOpts{})
		}

	}

	// our bundle is assembled, craft the final dir
	// (1) filters
	final := d.dag.Directory().WithDirectory("/", bundle, dagger.DirectoryWithDirectoryOpts{
		Include:   cfg.Include,
		Exclude:   cfg.Exclude,
		Gitignore: cfg.Gitignore,
		Owner:     cfg.Owner,
	})
	// (2) subpath selections
	final = final.Directory(cfg.TrimPrefix)

	if cfg.Patch != "" {
		final = final.WithPatch(cfg.Patch)
	} else if cfg.PatchFile.Exists() {
		f, _, err := d.hashFile(cfg.PatchFile, noCache)
		if err != nil {
			return nil, nil, err
		}
		final = final.WithPatchFile(f)
	}

	// memoize
	idx.dir = final
	d.cat.Store(idx, idx)

	return idx.dir, idx.cfg, nil
}

type exportImageFileConfig struct {
	Kind  string    `json:"$kind"`
	Name  string    `json:"name"`
	Path  string    `json:"path"`
	Tags  []string  `json:"tags"`
	Image cue.Value `json:"image"`
}

type exportImageFileIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *exportImageFileConfig
	ctr  *dagger.Container
}

func (idx *exportImageFileIndex) Key() string {
	if idx.cfg == nil {
		return "#exportImageFile.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#exportImageFile.%s", mk)
	}
	return fmt.Sprintf("#exportImageFile.%s", idx.cfg.Name)
}

func (d *Dag) HashExportImageFile(step cue.Value, noCache bool) (*dagger.Container, *exportImageFileConfig, error) {
	d.mx.RLock()
	var cfg exportImageFileConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, err
	}

	// index for query and create if not found
	idx := &exportImageFileIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*exportImageFileIndex)
		return ix.ctr, ix.cfg, nil
	}
	c, err := d.HashContainer(cfg.Image, noCache)
	if err != nil {
		return nil, nil, err
	}
	idx.ctr = c

	// memoize
	d.cat.Store(idx, idx)

	return idx.ctr, idx.cfg, nil
}

type exportImageConfig struct {
	Kind  string    `json:"$kind"`
	Name  string    `json:"name"`
	Reg   string    `json:"reg"`
	Tags  []string  `json:"tags"`
	Image cue.Value `json:"image"`
}

type exportImageIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *exportImageConfig
	ctr  *dagger.Container
}

func (idx *exportImageIndex) Key() string {
	if idx.cfg == nil {
		return "#exportImage.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#exportImage.%s", mk)
	}
	return fmt.Sprintf("#exportImage.%s", idx.cfg.Name)
}

func (d *Dag) HashExportImage(step cue.Value, noCache bool) (*dagger.Container, *exportImageConfig, error) {
	d.mx.RLock()
	var cfg exportImageConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, err
	}

	// index for query and create if not found
	idx := &exportImageIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*exportImageIndex)
		return ix.ctr, ix.cfg, nil
	}
	c, err := d.HashContainer(cfg.Image, noCache)
	if err != nil {
		return nil, nil, err
	}
	idx.ctr = c

	// memoize
	d.cat.Store(idx, idx)

	return idx.ctr, idx.cfg, nil
}

type publishImageConfig struct {
	Kind  string    `json:"$kind"`
	Name  string    `json:"name"`
	Reg   string    `json:"reg"`
	Tags  []string  `json:"tags"`
	Image cue.Value `json:"image"`
}

type publishImageIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *publishImageConfig
	ctr  *dagger.Container
}

func (idx *publishImageIndex) Key() string {
	if idx.cfg == nil {
		return "#publishImage.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#publishImage.%s", mk)
	}
	return fmt.Sprintf("#publishImage.%s", idx.cfg.Name)
}

func (d *Dag) HashPublishImage(step cue.Value, noCache bool) (*dagger.Container, *publishImageConfig, error) {
	d.mx.RLock()
	var cfg publishImageConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, err
	}

	// index for query and create if not found
	idx := &publishImageIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*publishImageIndex)
		return ix.ctr, ix.cfg, nil
	}
	c, err := d.HashContainer(cfg.Image, noCache)
	if err != nil {
		return nil, nil, err
	}
	idx.ctr = c

	// memoize
	d.cat.Store(idx, idx)

	return idx.ctr, idx.cfg, nil
}
