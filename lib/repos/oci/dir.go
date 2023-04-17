package oci

import (
	ignore "github.com/sabhiram/go-gitignore"
)

func NewDir(path string, ignores []string) Dir {
	var ign *ignore.GitIgnore
	if len(ignores) > 0 {
		ign = ignore.CompileIgnoreLines(ignores...)
	}

	return Dir{
		Path: path,
		ign:  ign,
	}
}

type Dir struct {
	ign  *ignore.GitIgnore
	Path string
}

func (d Dir) Excluded(rel string) bool {
	if d.ign == nil {
		return false
	}

	return d.ign.MatchesPath(rel)
}
