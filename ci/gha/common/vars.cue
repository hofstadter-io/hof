package common

Formatters: [
	"prettier",
	"csharpier",
	"black",
]

Versions: {
	docker: "20.x" | ["20.x", "23.x"]
	go: "1.20.x" | ["1.19.x", "1.20.x"]
	os: "ubuntu-latest" | ["ubuntu-latest", "macos-latest"]
}
