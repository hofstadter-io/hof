package oci

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/v1/types"
	ignore "github.com/sabhiram/go-gitignore"
)

const (
	modIgnoreFile = ".hofmod-ignore"
)

func NewDir(mediaType types.MediaType, relPath string, ignores []string) Dir {
	var ign *ignore.GitIgnore
	if len(ignores) > 0 {
		ign = ignore.CompileIgnoreLines(ignores...)
	}

	return Dir{
		mediaType: mediaType,
		relPath:   relPath,
		ign:       ign,
	}
}

type Dir struct {
	ign       *ignore.GitIgnore
	relPath   string
	mediaType types.MediaType
}

func (d Dir) Excluded(rel string) bool {
	if d.ign == nil {
		return false
	}

	return d.ign.MatchesPath(rel)
}

func NewDeps() Dir {
	return NewDir(HofstadterModuleDeps, "cue.mod", []string{
		"*",
		"!module.cue",
		"!sums.cue",
	})
}

func NewCode(workingDir string) (Dir, error) {
	ignores := []string{
		"cue.mod/pkg",
		".git",
	}

	p := filepath.Join(workingDir, modIgnoreFile)

	if _, err := os.Stat(p); err == nil {
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return Dir{}, fmt.Errorf("read file %s: %w", modIgnoreFile, err)
		}

		ls := strings.Split(string(b), "\n")
		ignores = append(ignores, ls...)
	}

	return NewDir(HofstadterModuleCode, "", ignores), nil
}
