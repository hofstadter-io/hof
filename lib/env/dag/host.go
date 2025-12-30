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
	return fmt.Sprintf("#hostFile.%s", idx.cfg.Name)
}

func (d *Dag) hashHostFile(step cue.Value) (*dagger.File, string, error) {
	var cfg hostFileConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, "", fmt.Errorf("while decoding hashHostFile: %w", err)
	}

	// index for query and create if not found
	idx := &hostFileIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostFileIndex)
		return ix.file, cfg.Path, nil
	}

	// load for realz
	idx.file = d.dag.Host().File(cfg.Path, dagger.HostFileOpts{
		NoCache: cfg.NoCache,
	})

	// memoize
	d.cat[idx] = idx

	return idx.file, cfg.Path, nil
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
	return fmt.Sprintf("#hostDir.%s", idx.cfg.Name)
}

func (d *Dag) hashHostDir(step cue.Value) (*dagger.Directory, string, error) {
	var cfg hostDirConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, "", fmt.Errorf("while decoding hashHostDir: %w", err)
	}

	// index for query and create if not found
	idx := &hostDirIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostDirIndex)
		return ix.dir, cfg.Path, nil
	}

	// load for realz
	idx.dir = d.dag.Host().Directory(cfg.Path, dagger.HostDirectoryOpts{
		// Include:   cfg.Include,
		// Exclude:   cfg.Exclude,
		// NoCache:   cfg.NoCache,
		Gitignore: cfg.GitIgnore,
	})

	// memoize
	d.cat[idx] = idx

	return idx.dir, cfg.Path, nil
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
	return fmt.Sprintf("#hostService.%s", idx.cfg.Name)
}

func (d *Dag) hashHostService(step cue.Value) (*dagger.Service, error) {
	var cfg hostServiceConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashHostService: %w", err)
	}

	// index for query and create if not found
	idx := &hostServiceIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostServiceIndex)
		return ix.svc, nil
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

	return idx.svc, nil
}

type hostTunnelConfig struct {
	Kind    string        `json:"$kind"`
	Name    string        `json:"name"`
	Service cue.Value     `json:"service"`
	Native  bool          `json:"native"`
	Ports   []portForward `json:"ports"`
}

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
	return fmt.Sprintf("#hostTunnel.%s", idx.cfg.Name)
}

func (d *Dag) hashHostTunnel(step cue.Value) (*dagger.Service, error) {
	var cfg hostTunnelConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashHostTunnel: %w", err)
	}

	// index for query and create if not found
	idx := &hostTunnelIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostTunnelIndex)
		return ix.svc, nil
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

	svc, _, err := d.hashService(cfg.Service)
	if err != nil {
		return nil, err
	}

	idx.svc = d.dag.Host().Tunnel(svc, dagger.HostTunnelOpts{
		Native: cfg.Native,
		Ports:  ports,
	})

	// memoize
	d.cat[idx] = idx

	return idx.svc, nil
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
	return fmt.Sprintf("#hostSocket.%s", idx.cfg.Name)
}

func (d *Dag) hashHostSocket(step cue.Value) (*dagger.Socket, error) {
	var cfg hostSocketConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashHostSocket: %w", err)
	}

	// index for query and create if not found
	idx := &hostSocketIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostSocketIndex)
		return ix.sock, nil
	}

	// load for realz
	idx.sock = d.dag.Host().UnixSocket(cfg.Path)

	// memoize
	d.cat[idx] = idx

	return idx.sock, nil
}

type portForward struct {
	Name     string `json:"name"`
	Backend  int    `json:"backend"`
	Frontend int    `json:"frontend"`
	Protocol string `json:"protocol"`
}
