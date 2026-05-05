package lang

flags: {
	goos: *"linux" | "darwin"
	arch: *"arm64" | "amd64"
	narch: [
		if arch == "amd64" {"x86"},
		arch,
	][0]

}