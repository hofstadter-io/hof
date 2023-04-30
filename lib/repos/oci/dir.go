package oci

import (
	"github.com/google/go-containerregistry/pkg/v1/types"
	ignore "github.com/sabhiram/go-gitignore"
)

func NewDir(mediaType types.MediaType, path string, ignores []string) Dir {
	var ign *ignore.GitIgnore
	if len(ignores) > 0 {
		ign = ignore.CompileIgnoreLines(ignores...)
	}

	return Dir{
		mediaType: mediaType,
		path:      path,
		ign:       ign,
	}
}

type Dir struct {
	ign       *ignore.GitIgnore
	path      string
	mediaType types.MediaType
}

func (d Dir) Excluded(rel string) bool {
	if d.ign == nil {
		return false
	}

	return d.ign.MatchesPath(rel)
}
