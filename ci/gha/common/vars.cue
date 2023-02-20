package common

Formatters: [
	"prettier",
	"csharpier",
	"black",
]

GoStrategy: {
	"fail-fast": false
	matrix: {
		"go-version": ["1.18.x", "1.19.x", "1.20.x"]
		os: ["ubuntu-latest", "macos-latest"]
	}
}
