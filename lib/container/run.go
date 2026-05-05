package container

import (
	"context"
	"fmt"
)

type Params struct {
	Name    Name `json:"name,omitempty"`
	Replace bool `json:"replace,omitempty"`

	Entrypoint    string   `json:"entrypoint,omitempty"`
	Env           []string `json:"env,omitempty"`
	Hostname      string   `json:"hostname,omitempty"`
	Mount         []string `json:"mount,omitempty"`
	Network       string   `json:"network,omitempty"`
	NetworkAlias  []string `json:"network-alias,omitempty"`
	NoHealthcheck bool     `json:"noHealthcheck,omitempty"`
	Platform      string   `json:"platform,omitempty"`
	Privileged    bool     `json:"privileged,omitempty"`
	Publish       []string `json:"publish,omitempty"`
	Pull          string   `json:"pull,omitempty"`
	Restart       string   `json:"restart,omitempty"`
	User          string   `json:"user,omitempty"`
	Volume        []string `json:"volume,omitempty"`
	Workdir       string   `json:"workdir,omitempty"`
	AddHost       []string `json:"addHost,omitempty"`

	Args []string `json:"args,omitempty"`
}

func (r runtime) Run(ctx context.Context, ref Ref, p *Params) error {
	if p.Replace {
		if err := r.Remove(ctx, p.Name); err != nil {
			return fmt.Errorf("remove: %w", err)
		}
	}

	args := []string{"run", "--detach"}

	if p.Name != "" {
		args = append(args, "--name", string(p.Name))
	}

	if p.Entrypoint != "" {
		args = append(args, "--entrypoint", p.Entrypoint)
	}

	for _, e := range p.Env {
		args = append(args, "--env", e)
	}

	if p.Hostname != "" {
		args = append(args, "--hostname", p.Hostname)
	}

	for _, m := range p.Mount {
		args = append(args, "--mount", m)
	}

	if p.Network != "" {
		args = append(args, "--network", p.Network)
	}

	for _, na := range p.NetworkAlias {
		args = append(args, "--network-alias", na)
	}

	if p.NoHealthcheck {
		args = append(args, "--no-healthcheck")
	}

	if p.Platform != "" {
		args = append(args, "--platform", p.Platform)
	}

	if p.Privileged {
		args = append(args, "--privileged")
	}

	for _, pblsh := range p.Publish {
		args = append(args, "--publish", pblsh)
	}

	if p.Pull != "" {
		args = append(args, "--pull", p.Pull)
	}

	if p.Restart != "" {
		args = append(args, "--restart", p.Restart)
	}

	if p.User != "" {
		args = append(args, "--user", p.User)
	}

	for _, v := range p.Volume {
		args = append(args, "--volume", v)
	}

	if p.Workdir != "" {
		args = append(args, "--workdir", p.Workdir)
	}

	for _, eh := range p.AddHost {
		args = append(args, "--add-host", eh)
	}

	args = append(args, string(ref))
	args = append(args, p.Args...)

	if _, err := r.exec(ctx, args...); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
