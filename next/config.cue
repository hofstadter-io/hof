package docs

import "strings"

version: "v0.6.9-beta.1"

// download links
_url: "https://github.com/hofstadter-io/hof/releases/download/\(version)/hof_\(version)"
_matrix: {
	os: ["darwin", "linux", "windows"]
	arch: ["amd64", "arm64"]
}
downloads: [
	for _, _os in _matrix.os for _, _arch in _matrix.arch {
		os:    _os
		arch:  _arch
		key:   "\(_os)_\(_arch)"
		title: strings.ToTitle("\(_os) \(_arch)")
		url:   "\(_url)_\(key)"
	},
]
