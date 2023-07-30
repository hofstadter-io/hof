package container

import (
	"context"
	"fmt"
	"time"
)

type RuntimeVersion struct {
	Name string
	Client struct {
		Version    string
		APIVersion string
	}
	Server struct {
		Version       string
		APIVersion    string
		MinAPIVersion string
	}
}

func GetVersion() (RuntimeVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return rt.Version(ctx)
}

func (r runtime) Version(ctx context.Context) (RuntimeVersion, error) {
	var rv RuntimeVersion
	if err := r.execJSON(ctx, &rv, "version", "--format", "{{ json . }}"); err != nil {
		return rv, fmt.Errorf("exec json: %w", err)
	}

	rv.Name = string(r.bin)
	return rv, nil
}

func (r RuntimeVersion) String() string {
	return fmt.Sprintf(
		"%s [%s (client) | %s (server)]",
		r.Name,
		r.Client.Version,
		r.Server.Version,
	)
}
