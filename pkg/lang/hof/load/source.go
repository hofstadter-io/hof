// Copyright 2019 CUE Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package load

import (
	"github.com/hofstadter-io/hof/pkg/lang/hof/ast"
	"github.com/hofstadter-io/hof/pkg/lang/hof/format"
)

// A Source represents file contents.
type Source interface {
	contents() ([]byte, error)
}

// FromString creates a Source from the given string.
func FromString(s string) Source {
	return stringSource(s)
}

// FromBytes creates a Source from the given bytes. The contents are not
// copied and should not be modified.
func FromBytes(b []byte) Source {
	return bytesSource(b)
}

// FromFile creates a Source from the given *ast.File. The file should not be
// modified. It is assumed the file is error-free.
func FromFile(f *ast.File) Source {
	return (*fileSource)(f)
}

type stringSource string

func (s stringSource) contents() ([]byte, error) {
	return []byte(s), nil
}

type bytesSource []byte

func (s bytesSource) contents() ([]byte, error) {
	return []byte(s), nil
}

type fileSource ast.File

func (s *fileSource) contents() ([]byte, error) {
	return format.Node((*ast.File)(s))
}
