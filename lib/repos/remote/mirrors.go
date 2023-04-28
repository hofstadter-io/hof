package remote

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hofstadter-io/hof/lib/repos/git"
	"github.com/hofstadter-io/hof/lib/repos/oci"
)

var (
	mirrorsGit = []string{
		"github.com",
		"gitlab.com",
		"bitbucket.com",
	}
	mirrorsOCI = []string{}
)

const (
	hofDir             = "hof"
	mirrorsFileName    = "mirrors.json"
	mirrorsFileNameEnv = "HOF_MOD_MIRRORFILE"
)

func mirrorsFilePath() (string, error) {
	p := os.Getenv(mirrorsFileNameEnv)
	if p != "" {
		return p, nil
	}

	d, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user cache dir: %w", err)
	}

	return filepath.Join(d, hofDir, mirrorsFileName), nil
}

func NewMirrors() (*Mirrors, error) {
	p, err := mirrorsFilePath()
	if err != nil {
		return nil, fmt.Errorf("mirrors file path: %w", err)
	}

	if _, err = os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return &Mirrors{
			values: make(map[Kind][]string),
		}, nil
	}

	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("os open %s: %w", p, err)
	}

	defer f.Close()

	var m Mirrors
	if err := json.NewDecoder(f).Decode(&m.values); err != nil {
		return nil, fmt.Errorf("json decode %s: %w", p, err)
	}

	return &m, nil
}

type Mirrors struct {
	valuesMu sync.RWMutex
	values   map[Kind][]string
}

func (m *Mirrors) Is(ctx context.Context, k Kind, mod string) (bool, error) {
	var (
		mirrors  []string
		netCheck func(context.Context, string) (bool, error)
	)

	switch k {
	case KindGit:
		mirrors = mirrorsGit
		netCheck = m.netCheckGit
	case KindOCI:
		mirrors = mirrorsOCI
		netCheck = m.netCheckOCI
	default:
		return false, fmt.Errorf("unknow kind: %s", k)
	}

	for _, ss := range mirrors {
		if mod == ss {
			return true, nil
		}
	}

	if m.hasValue(k, mod) {
		return true, nil
	}

	is, err := netCheck(ctx, mod)
	if err != nil {
		return false, err
	}

	if is {
		m.valuesMu.Lock()
		m.values[k] = append(m.values[k], mod)
		m.valuesMu.Unlock()
	}

	return is, nil
}

func (m *Mirrors) hasValue(k Kind, mod string) bool {
	m.valuesMu.RLock()
	defer m.valuesMu.RUnlock()

	if vals, ok := m.values[k]; ok {
		for _, v := range vals {
			if mod == v {
				return true
			}
		}
	}

	return false
}

func (m *Mirrors) netCheckGit(ctx context.Context, mod string) (bool, error) {
	return git.IsNetworkReachable(ctx, mod)
}

func (m *Mirrors) netCheckOCI(ctx context.Context, mod string) (bool, error) {
	return oci.IsNetworkReachable(mod), nil
}

func (m *Mirrors) Set(k Kind, s string) {
	m.valuesMu.Lock()
	defer m.valuesMu.Unlock()

	if vals, ok := m.values[k]; ok {
		vals = append(vals, s)
		m.values[k] = vals
	}
}

func (m *Mirrors) Close() error {
	m.valuesMu.Lock()
	defer m.valuesMu.Unlock()

	if len(m.values) == 0 {
		return nil
	}

	p, err := mirrorsFilePath()
	if err != nil {
		return fmt.Errorf("mirrors file path: %w", err)
	}

	dir := filepath.Dir(p)
	if err = os.MkdirAll(dir, 0o600); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("os open file %s: %w", p, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&m.values); err != nil {
		return fmt.Errorf("json encode %s: %w", p, err)
	}

	return nil
}
