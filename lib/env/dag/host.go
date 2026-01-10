package dag

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hostImageConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`
}

type hostImageIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hostImageConfig
	img  *dagger.Container
}

func (idx *hostImageIndex) Key() string {
	if idx.cfg == nil {
		return "#hostImage.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#hostImage.%s", mk)
	}
	return fmt.Sprintf("#hostImage.%s", idx.cfg.Name)
}

func (d *Dag) HashHostImage(step cue.Value, noCache bool) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg hostImageConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding hashHostImage: %w", err)
	}

	// index for query and create if not found
	idx := &hostImageIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hostImageIndex)
		return ix.img, nil
	}

	// load for realz
	idx.img = d.dag.Host().ContainerImage(cfg.Name)

	// memoize
	d.cat.Store(idx, idx)

	return idx.img, nil
}

type hostFileConfig struct {
	Kind    string `json:"$kind"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	NoCache bool   `json:"noCache"`
}

type hostFileIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hostFileConfig
	file *dagger.File
}

func (idx *hostFileIndex) Key() string {
	if idx.cfg == nil {
		return "#hostFile.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#hostFile.%s", mk)
	}
	return fmt.Sprintf("#hostFile.%s", idx.cfg.Name)
}

func (d *Dag) HashHostFile(val cue.Value, noCache bool) (*dagger.File, *hostFileConfig, error) {
	d.mx.RLock()
	var cfg hostFileConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostFile: %w", err)
	}

	// index for query and create if not found
	idx := &hostFileIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hostFileIndex)
		return ix.file, ix.cfg, nil
	}

	// load for realz
	idx.file = d.dag.Host().File(cfg.Path, dagger.HostFileOpts{
		NoCache: cfg.NoCache || noCache,
	})

	// memoize
	d.cat.Store(idx, idx)

	return idx.file, idx.cfg, nil
}

type hostDirConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`
	Path string `json:"path"`

	TrimPrefix string    `json:"trimPrefix"`
	Patch      string    `json:"patch"`
	PatchFile  cue.Value `json:"patchFile"`

	Include   []string `json:"include"`
	Exclude   []string `json:"exclude"`
	GitIgnore bool     `json:"gitignore"`
	NoCache   bool     `json:"noCache"`
}

type hostDirIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hostDirConfig
	dir  *dagger.Directory
}

func (idx *hostDirIndex) Key() string {
	if idx.cfg == nil {
		return "#hostDir.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#hostDir.%s", mk)
	}
	return fmt.Sprintf("#hostDir.%s", idx.cfg.Name)
}

func (d *Dag) HashHostDir(val cue.Value, noCache bool) (*dagger.Directory, *hostDirConfig, error) {
	d.mx.RLock()
	var cfg hostDirConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostDir: %w", err)
	}

	// index for query and create if not found
	idx := &hostDirIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hostDirIndex)
		return ix.dir, ix.cfg, nil
	}

	// load for realz
	final := d.dag.Host().Directory(cfg.Path, dagger.HostDirectoryOpts{
		Include:   cfg.Include,
		Exclude:   cfg.Exclude,
		NoCache:   cfg.NoCache || noCache,
		Gitignore: cfg.GitIgnore,
	})

	// (2) subpath selections
	if cfg.TrimPrefix != "" {
		final = final.Directory(cfg.TrimPrefix)
	}

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

type hostServiceConfig struct {
	Kind  string        `json:"$kind"`
	Name  string        `json:"name"`
	Host  string        `json:"host"`
	Ports []portForward `json:"ports"`
}

type hostServiceIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hostServiceConfig
	svc  *dagger.Service
}

func (idx *hostServiceIndex) Key() string {
	if idx.cfg == nil {
		return "#hostService.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#hostService.%s", mk)
	}
	return fmt.Sprintf("#hostService.%s", idx.cfg.Name)
}

func (d *Dag) HashHostService(val cue.Value) (*dagger.Service, *hostServiceConfig, error) {
	d.mx.RLock()
	var cfg hostServiceConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostService: %w", err)
	}

	// index for query and create if not found
	idx := &hostServiceIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hostServiceIndex)
		return ix.svc, ix.cfg, nil
	}

	// load for realz
	ports := []dagger.PortForward{}
	for _, p := range cfg.Ports {
		ports = append(ports, dagger.PortForward{
			Protocol: dagger.NetworkProtocol(strings.ToUpper(p.Protocol)),
			Frontend: p.Frontend,
			Backend:  p.Backend,
		})
	}
	idx.svc = d.dag.Host().Service(ports, dagger.HostServiceOpts{
		Host: cfg.Host,
	})

	// memoize
	d.cat.Store(idx, idx)

	return idx.svc, idx.cfg, nil
}

type hostTunnelConfig struct {
	Kind    string        `json:"$kind"`
	Name    string        `json:"name"`
	Service cue.Value     `json:"service"`
	Native  bool          `json:"native"`
	Ports   []portForward `json:"ports"`
}

// maybe we should just export and return this, things are getting more mature
type hostTunnelIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hostTunnelConfig
	svc  *dagger.Service
}

func (idx *hostTunnelIndex) Key() string {
	if idx.cfg == nil {
		return "#hostTunnel.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#hostTunnel.%s", mk)
	}
	return fmt.Sprintf("#hostTunnel.%s", idx.cfg.Name)
}

func (d *Dag) HashHostTunnel(val cue.Value) (*dagger.Service, *hostTunnelConfig, error) {
	d.mx.RLock()
	var cfg hostTunnelConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostTunnel: %w", err)
	}

	// index for query and create if not found
	idx := &hostTunnelIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hostTunnelIndex)
		return ix.svc, ix.cfg, nil
	}

	// load for realz
	ports := []dagger.PortForward{}
	for _, p := range cfg.Ports {
		ports = append(ports, dagger.PortForward{
			Protocol: dagger.NetworkProtocol(strings.ToUpper(p.Protocol)),
			Frontend: p.Frontend,
			Backend:  p.Backend,
		})
	}

	svc, _, err := d.HashService(cfg.Service, false) // hmm, wonder why it ai chose false here and didn't pass it in, host interaction is a weird one and really shouldn't cache every in these networky ones anyway?
	if err != nil {
		return nil, nil, err
	}

	idx.svc = d.dag.Host().Tunnel(svc, dagger.HostTunnelOpts{
		Native: cfg.Native,
		Ports:  ports,
	})

	// memoize
	d.cat.Store(idx, idx)

	return idx.svc, idx.cfg, nil
}

type hostSocketConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type hostSocketIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hostSocketConfig
	sock *dagger.Socket
}

func (idx *hostSocketIndex) Key() string {
	if idx.cfg == nil {
		return "#hostSocket.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#hostSocket.%s", mk)
	}
	return fmt.Sprintf("#hostSocket.%s", idx.cfg.Name)
}

func (d *Dag) HashHostSocket(val cue.Value) (*dagger.Socket, *hostSocketConfig, error) {
	d.mx.RLock()
	var cfg hostSocketConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostSocket: %w", err)
	}

	// index for query and create if not found
	idx := &hostSocketIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hostSocketIndex)
		return ix.sock, ix.cfg, nil
	}

	// load for realz
	idx.sock = d.dag.Host().UnixSocket(cfg.Path)

	// memoize
	d.cat.Store(idx, idx)

	return idx.sock, idx.cfg, nil
}

type portForward struct {
	Name     string `json:"name"`
	Backend  int    `json:"backend"`
	Frontend int    `json:"frontend"`
	Protocol string `json:"protocol"`
}

type stepUnixSocketConfig struct {
	Kind string `json:"$kind"`

	// args
	Path   string    `json:"path"`
	Source cue.Value `json:"source"`

	// opts (depending on source type?)
	Owner  string `json:"owner"`
	Expand bool   `json:"expand"`
}

func (d *Dag) stepUnixSocketHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {

	// DEV HACK
	// return c, nil

	d.mx.RLock()
	var cfg stepUnixSocketConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepUnixSocketConfig: %w", err)
	}
	k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return c, fmt.Errorf("missing $kind in stepUnixSocketConfig.source: %v, got %v", step, k)
	}

	var sock *dagger.Socket
	ks, _ := k.String()
	switch ks {
	case "#hostSocket":
		sock, _, err = d.HashHostSocket(cfg.Source)

	default:
		return c, fmt.Errorf("unsupported $kind in stepUnixSocketConfig.source: %v", step)
	}

	if err != nil {
		return nil, err
	}

	c = c.WithUnixSocket(cfg.Path, sock, dagger.ContainerWithUnixSocketOpts{
		Owner:  cfg.Owner,
		Expand: cfg.Expand,
	})

	return c, nil
}

type hostExecConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`

	Args    []string `json:"args"`
	Workdir string   `json:"workdir"`
	Envs    []string `json:"envs"`
	Stdin   string   `json:"stdin"`
	Stdout  string   `json:"stdout"`
	Stderr  string   `json:"stderr"`
	AllEnv  bool     `json:"allEnv"`
}

func (d *Dag) HashHostExec(val cue.Value) error {
	d.mx.RLock()
	var cfg hostExecConfig
	err := val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return fmt.Errorf("while decoding hashHostExec: %w", err)
	}

	if len(cfg.Args) == 0 {
		return fmt.Errorf("hostExec: args must have at least one element")
	}

	cmd := exec.Command(cfg.Args[0], cfg.Args[1:]...)
	cmd.Dir = cfg.Workdir
	cmd.Env = cfg.Envs
	if cfg.AllEnv {
		cmd.Env = append(os.Environ(), cmd.Env...)
	}

	if cfg.Stdin != "" {
		f, err := os.Open(cfg.Stdin)
		if err != nil {
			return err
		}
		defer f.Close()
		cmd.Stdin = f
	}
	if cfg.Stdout != "" {
		f, err := os.Create(cfg.Stdout)
		if err != nil {
			return err
		}
		defer f.Close()
		cmd.Stdout = f
	}
	if cfg.Stderr != "" {
		f, err := os.Create(cfg.Stderr)
		if err != nil {
			return err
		}
		defer f.Close()
		cmd.Stderr = f
	}

	return cmd.Run()
}
