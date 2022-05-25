// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sumfile

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type Sum struct {
	Mods map[Version][]string
}

type Version struct {
	Path    string
	Version string
}

// emptyGoModHash is the hash of a 1-file tree containing a 0-length go.mod.
// A bug caused us to write these into go.sum files for non-modules.
// We detect and remove them.
const emptyGoModHash = "h1:G7mAYYxgmS0lVkHyy2hEOLQCFB0DlQFTMLWggykrydY="

func ParseSum(data []byte, file string) (Sum, error) {
	var sum Sum
	sum.Mods = make(map[Version][]string)

	lineno := 0
	for len(data) > 0 {
		var line []byte
		lineno++
		i := bytes.IndexByte(data, '\n')
		if i < 0 {
			line, data = data, nil
		} else {
			line, data = data[:i], data[i+1:]
		}
		f := strings.Fields(string(line))
		if len(f) == 0 {
			// blank line; skip it
			continue
		}
		if len(f) != 3 {
			return sum, fmt.Errorf("malformed %s:\n%s:%d: wrong number of fields %v", file, file, lineno, len(f))
		}
		if f[2] == emptyGoModHash {
			// Old bug; drop it.
			continue
		}
		mod := Version{Path: f[0], Version: f[1]}
		sum.Mods[mod] = append(sum.Mods[mod], f[2])
	}

	return sum, nil
}

func (sum *Sum) Add(ver Version, hash string) {
	if sum.Mods == nil {
		sum.Mods = make(map[Version][]string)
	}
	val, ok := sum.Mods[ver]
	if !ok {
		val = make([]string, 0)
	}
	val = append(val, hash)
	sum.Mods[ver] = val
}

func (sum *Sum) Write() (string, error) {
	var w strings.Builder
	// build up slice
	var sorted []Version
	for ver, _ := range sum.Mods {
		sorted = append(sorted, ver)
	}

	// sort slice by ver.Path
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Path == sorted[j].Path {
			return sorted[i].Version < sorted[j].Version
		}
		return sorted[i].Path < sorted[j].Path
	})

	// print
	for _, ver := range sorted {
		list := sum.Mods[ver]
		fmt.Fprintln(&w, ver.Path, ver.Version, list[0])
	}

	return w.String(), nil
}
