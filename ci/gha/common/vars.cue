package common

Formatters: [
	"prettier",
	"csharpier",
	"black",
]

Versions: {
	docker: "24.0.7"
	go:     "1.23.x" | ["1.23.x", "1.22.x"]
	os:     "ubuntu-latest" | ["ubuntu-latest"]
}
