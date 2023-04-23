package remote

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
)

var (
	knownGit = []string{
		"github.com",
		"gitlab.com",
		"bitbucket.com",
	}
	knownOCI = []string{}
)

const knownsFileName = "hof-knowns.json"

func knownsFilePath() (string, error) {
	d, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("user cache dir: %w", err)
	}

	return path.Join(d, knownsFileName), nil
}

func NewKnowns() (*Knowns, error) {
	p, err := knownsFilePath()
	if err != nil {
		return nil, fmt.Errorf("knowns file path: %w", err)
	}

	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("os open %s: %w", p, err)
	}

	defer f.Close()

	var kn Knowns
	if err := json.NewDecoder(f).Decode(&kn.values); err != nil {
		return nil, fmt.Errorf("json decode %s: %w", p, err)
	}

	return &kn, nil
}

type Knowns struct {
	valuesMu sync.RWMutex
	values   map[Kind][]string
}

func (kn *Knowns) IsKnown(k Kind, s string) bool {
	var builtins []string
	switch k {
	case KindGit:
		builtins = knownGit
	case KindOCI:
		builtins = knownOCI
	}

	for _, ss := range builtins {
		if s == ss {
			return true
		}
	}

	kn.valuesMu.RLock()
	defer kn.valuesMu.RUnlock()

	if vals, ok := kn.values[k]; ok {
		for _, v := range vals {
			if s == v {
				return true
			}
		}
	}

	return false
}

func (kn *Knowns) Set(k Kind, s string) {
	kn.valuesMu.Lock()
	defer kn.valuesMu.Unlock()

	if vals, ok := kn.values[k]; ok {
		vals = append(vals, s)
		kn.values[k] = vals
	}
}

func (kn *Knowns) Close() error {
	kn.valuesMu.Lock()
	defer kn.valuesMu.Unlock()

	p, err := knownsFilePath()
	if err != nil {
		return fmt.Errorf("knowns file path: %w", err)
	}

	f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("os open file %s: %w", p, err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(&kn.values); err != nil {
		return fmt.Errorf("json encode %s: %w", p, err)
	}

	return nil
}
