// Copyright 2020 Hofstadter, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mappingfile

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type Mappings struct {
	Mods map[Module][]string
}

type Module struct {
	Explicit bool
	Path     string
	Version  string
}

func ParseMapping(data []byte, file string) (Mappings, error) {
	var mappings Mappings
	mappings.Mods = make(map[Module][]string)

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
			return mappings, fmt.Errorf("malformed %s:\n%s:%d: wrong number of fields %v", file, file, lineno, len(f))
		}

		mod := Module{Path: f[0], Version: f[1]}
		mappings.Mods[mod] = append(mappings.Mods[mod], f[2])
	}

	return mappings, nil
}

func (mappings *Mappings) Print() error {
	// build up slice
	var sorted []Module
	for ver, _ := range mappings.Mods {
		sorted = append(sorted, ver)
	}

	// sort slice by ver.Path
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	// print
	for _, ver := range sorted {
		list := mappings.Mods[ver]
		fmt.Println(ver.Path, ver.Version, ver.Explicit, list)
	}

	return nil
}
