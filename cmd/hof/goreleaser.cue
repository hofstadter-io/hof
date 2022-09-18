project_name: "hof"

release: {
	disable: false
	draft:   false
	github: {
		name:  "hof"
		owner: "hofstadter-io"
	}
}

changelog: {
	filters: exclude: ["^docs:", "^test:"]
	sort: "asc"
}

checksum: name_template: "{{ .ProjectName }}_{{ .Version }}_checksums.txt"
snapshot: name_template: "{{ .Tag }}-SNAPSHOT-{{ .ShortCommit }}"

builds: [{
	binary: "hof"
	env: ["CGO_ENABLED=0"]
	goarch: ["amd64", "arm64"]
	goos: ["darwin", "linux", "windows"]
	_flags: [
		"Version={{ .Version }}",
		"Commit={{ .FullCommit }}",
		"BuildDate={{ .Date }}",
		"BuildOS={{ .Os }}",
		"BuildArch={{ .Arch }}",
		"BuildArm={{ .Arm }}",
	]
	ldflags: [ "-s -w", for f in _flags { "-X github.com/hofstadter-io/hof/cmd/hof/verinfo.\(f)" } ]
	main: "main.go"
}]

dockers: [...{
	skip_push: false
}] & [
	// hof images
	for cfg in [
		{base: "debian", suf: "" },
		{base: "debian", suf: "debian-" },
		{base: "alpine", suf: "alpine-" },
	] {
		dockerfile: "../../ci/hof/docker/Dockerfile.\(cfg.base)"
		image_templates: [ for suf in ["{{.Tag}}", "v{{ .Major }}.{{ .Minor }}", "{{ .ShortCommit }}", "latest"] {
			"hofstadter/hof:\(cfg.suf)\(suf)",
		}]
	}
]

archives: [{
	// this makes it so a binary only is uploaded, rather than a tar file
	files: ["thisfiledoesnotexist*"]
	format: "binary"
	replacements: {
		amd64:   "x86_64"
		darwin:  "Darwin"
		linux:   "Linux"
		windows: "Windows"
	}
}]

