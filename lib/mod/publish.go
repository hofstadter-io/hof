package mod

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/lib/repos/remote"
)

func Publish(taggedMod string) error {
	parts := strings.Split(taggedMod, ":")

	var (
		mod string
		tag string
	)

	switch {
	case len(parts) == 1:
		mod = taggedMod
		tag = "latest"
	case len(parts) == 2:
		mod = parts[0]
		tag = parts[1]
	default:
		return errors.New("invalid mod")
	}

	taggedMod = fmt.Sprintf("%s:%s", mod, tag)

	rmt, err := remote.Parse(mod)
	if err != nil {
		return fmt.Errorf("remote parse: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os get wd: %w", err)
	}

	d := filepath.Join(wd, "cue.mod", "pkg", mod)

	ctx := context.Background()
	if err = rmt.Publish(ctx, d, taggedMod); err != nil {
		return fmt.Errorf("remote publish: %w", err)
	}

	return nil
}
