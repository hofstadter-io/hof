package common

Formatters: [
	"prettier",
	"csharpier",
	"black",
]

Versions: {
	docker: "24.0.5"
	go:     "1.21.x" | ["1.20.x", "1.21.x"]
	os:     "ubuntu-latest" | ["ubuntu-latest", "macos-latest"]
}
