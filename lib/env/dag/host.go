package dag

import (
	"fmt"
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

func (d *Dag) HashHostImage(step cue.Value) (*dagger.Container, error) {
	var cfg hostImageConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashHostImage: %w", err)
	}

	// index for query and create if not found
	idx := &hostImageIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostImageIndex)
		return ix.img, nil
	}

	// load for realz
	idx.img = d.dag.Host().ContainerImage(cfg.Name)

	// memoize
	d.cat[idx] = idx

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

func (d *Dag) HashHostFile(val cue.Value) (*dagger.File, *hostFileConfig, error) {
	var cfg hostFileConfig
	err := val.Decode(&cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostFile: %w", err)
	}

	// index for query and create if not found
	idx := &hostFileIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostFileIndex)
		return ix.file, ix.cfg, nil
	}

	// load for realz
	idx.file = d.dag.Host().File(cfg.Path, dagger.HostFileOpts{
		NoCache: cfg.NoCache,
	})

	// memoize
	d.cat[idx] = idx

	return idx.file, idx.cfg, nil
}

type hostDirConfig struct {
	Kind      string   `json:"$kind"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	NoCache   bool     `json:"noCache"`
	Include   []string `json:"include"`
	Exclude   []string `json:"exclude"`
	GitIgnore bool     `json:"gitignore"`
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

func (d *Dag) HashHostDir(val cue.Value) (*dagger.Directory, *hostDirConfig, error) {
	var cfg hostDirConfig
	err := val.Decode(&cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostDir: %w", err)
	}

	// index for query and create if not found
	idx := &hostDirIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostDirIndex)
		return ix.dir, ix.cfg, nil
	}

	// load for realz
	idx.dir = d.dag.Host().Directory(cfg.Path, dagger.HostDirectoryOpts{
		Include:   cfg.Include,
		Exclude:   cfg.Exclude,
		NoCache:   cfg.NoCache,
		Gitignore: cfg.GitIgnore,
	})

	// memoize
	d.cat[idx] = idx

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
	var cfg hostServiceConfig
	err := val.Decode(&cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostService: %w", err)
	}

	// index for query and create if not found
	idx := &hostServiceIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
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
	d.cat[idx] = idx

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
	var cfg hostTunnelConfig
	err := val.Decode(&cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostTunnel: %w", err)
	}

	// index for query and create if not found
	idx := &hostTunnelIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
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

	svc, _, err := d.HashService(cfg.Service)
	if err != nil {
		return nil, nil, err
	}

	idx.svc = d.dag.Host().Tunnel(svc, dagger.HostTunnelOpts{
		Native: cfg.Native,
		Ports:  ports,
	})

	// memoize
	d.cat[idx] = idx

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
	var cfg hostSocketConfig
	err := val.Decode(&cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashHostSocket: %w", err)
	}

	// index for query and create if not found
	idx := &hostSocketIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostSocketIndex)
		return ix.sock, ix.cfg, nil
	}

	// load for realz
	idx.sock = d.dag.Host().UnixSocket(cfg.Path)

	// memoize
	d.cat[idx] = idx

	return idx.sock, idx.cfg, nil
}

type portForward struct {
	Name     string `json:"name"`
	Backend  int    `json:"backend"`
	Frontend int    `json:"frontend"`
	Protocol string `json:"protocol"`
}
