package dag

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/fsnotify/fsnotify"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashCacheConfig struct {
	Kind   string    `json:"$kind"`
	Name   string    `json:"name"`
	Source cue.Value `json:"source"`
	Watch  bool      `json:"watch"`
}

type hashCacheIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashCacheConfig
	vol  *dagger.CacheVolume
	dir  *dagger.Directory

	watch *hostDirConfig
	quit  chan bool
	ctr   *dagger.Container
}

func (idx *hashCacheIndex) Key() string {
	if idx.cfg == nil {
		return "#cache.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#cache.%s", mk)
	}
	return fmt.Sprintf("#cache.%s", idx.cfg.Name)
}

func (d *Dag) hashCache(step cue.Value) (*hashCacheIndex, error) {
	d.mx.RLock()
	var cfg hashCacheConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding hashCache: %w", err)
	}

	// index for query and create if not found
	idx := &hashCacheIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashCacheIndex)
		return ix, nil
	}

	// create our cache volume
	idx.vol = d.dag.CacheVolume(cfg.Name)

	k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if k.Exists() {
		var dir *dagger.Directory
		var doWatch bool // only possible with host dir (part1)
		ks, _ := k.String()
		switch ks {
		case "#dir":
			dir, _, err = d.hashDir(cfg.Source, false)
			if err != nil {
				return nil, err
			}

		case "#hostDir":
			var dcfg *hostDirConfig
			dir, dcfg, err = d.HashHostDir(cfg.Source, false)
			if err != nil {
				return nil, err
			}
			doWatch = cfg.Watch // we set to what the user set (part2)
			if doWatch {
				idx.watch = dcfg
			}

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

		}

		cfg.Watch = doWatch // this allows us to default to false and only allow true through on cases we want (part3)
		idx.dir = dir
		if doWatch {
			// debug
			fmt.Println("watch!", step)
			fmt.Println("watching:", idx.watch.Path, idx.watch.Include, idx.watch.Exclude, idx.watch.TrimPrefix)

			// list of paths to watch
			paths, err := d.gatherPaths(idx)
			if err != nil {
				return idx, err
			}
			fmt.Println("paths:", paths)

			// start the watcher
			idx.quit = make(chan bool)
			go func() {

				// setup container with cache
				c := d.dag.Container().From("alpine:latest")
				c = c.WithMountedCache(filepath.Join("/work", idx.watch.TrimPrefix), idx.vol, dagger.ContainerWithMountedCacheOpts{
					Source: idx.dir,
				})

				// our fs event handler
				handle := func(evt fsnotify.Event) error {
					fmt.Println("fsnotify:", evt.Op, evt.Name)
					// TODO, dir operations on c, need mutex here
					if evt.Op&fsnotify.Write == fsnotify.Write {
						path := evt.Name
						fmt.Println("updating:", path)
						bs, err := os.ReadFile(path)
						if err != nil {
							return err
						}
						_, _p := filepath.Split(path)
						f := d.dag.File(_p, string(bs))
						fmt.Println("...", path, _p, len(bs))
						tp, wp := filepath.Join("/tmp", path), filepath.Join("/work", path)
						// c, err = c.WithFile(tp, f).WithExec([]string{"sh", "-c", fmt.Sprintf("cp %s %s", tp, wp)}).Sync(d.ctx)
						c, err = c.WithFile(tp, f).Sync(d.ctx)
						if err != nil {
							return err
						}
						_, err = c.WithExec([]string{"sh", "-c", fmt.Sprintf("cp %s %s", tp, wp)}).Sync(d.ctx)
						if err != nil {
							return err
						}

					}
					return nil
				}

				err = startWatcher(handle, paths, time.Duration(0), idx.quit, false)
				if err != nil {
					fmt.Println("Error! while starting watcher:", err)
				}
			}()
			d.watched = append(d.watched, idx)
		}
	}

	// memoize
	d.cat.Store(idx, idx)

	return idx, nil
}

type stepTempConfig struct {
	Kind   string `json:"$kind"`
	Path   string `json:"path"`
	Size   int    `json:"size"`
	Expand bool   `json:"expand"`
}

func (d *Dag) stepTempHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepTempConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepTerm: %w", err)
	}

	c = c.WithMountedTemp(cfg.Path, dagger.ContainerWithMountedTempOpts{
		Size:   cfg.Size,
		Expand: cfg.Expand,
	})

	return c, nil
}

type stepMountConfig struct {
	Kind string `json:"$kind"`
	// args
	Path   string    `json:"path"`
	Source cue.Value `json:"source"`

	// opts (depending on source type?)
	Owner  string `json:"owner"`
	Expand bool   `json:"expand"`
	Mode   int    `json:"mode"`
}

func (d *Dag) stepMountHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepMountConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return c, err
	}

	// look for kind
	k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return c, fmt.Errorf("missing $kind in stepDir source: %v", step)
	}
	ks, _ := k.String()
	switch ks {
	case "#cache":
		idx, err := d.hashCache(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedCache(cfg.Path, idx.vol, dagger.ContainerWithMountedCacheOpts{
			Source: idx.dir,
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#secret":
		shh, err := d.hashSecret(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedSecret(cfg.Path, shh, dagger.ContainerWithMountedSecretOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
			Mode:   cfg.Mode,
		})

	case "#file":
		file, _, err := d.hashFile(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedFile(cfg.Path, file, dagger.ContainerWithMountedFileOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#hostFile":
		file, _, err := d.HashHostFile(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedFile(cfg.Path, file, dagger.ContainerWithMountedFileOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#dir":
		dir, _, err := d.hashDir(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#hostDir":
		dir, _, err := d.HashHostDir(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#gitRepo":
		repo, rcfg, rerr := d.hashGitRepo(cfg.Source, false)
		if rerr == nil {
			var dir *dagger.Directory
			if rcfg != nil && rcfg.Ref != "" {
				dir = repo.Ref(rcfg.Ref).Tree()
			} else {
				dir = repo.Head().Tree()
			}
			c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
				Owner:  cfg.Owner,
				Expand: cfg.Expand,
			})
		} else {
			err = rerr
		}
	case "#rootfs":
		dir, err := d.hashRootFS(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	default:
		return c, fmt.Errorf("unsupported $kind in stepMount.source: %v", step)
	}

	return c, nil
}
