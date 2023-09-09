project_name: "hof"

release: {
	disable: false
	draft:   true
	github: {
		name:  "hof"
		owner: "hofstadter-io"
	}
}

changelog: {
	filters: exclude: ["^docs:", "^test:"]
	sort: "asc"
}

checksum: name_template: "{{ .ProjectName }}_{{ .Tag }}_checksums.txt"
snapshot: name_template: "{{ .Tag }}-SNAPSHOT-{{ .ShortCommit }}"

builds: [{
	binary: "hof"
	env: ["CGO_ENABLED=0"]
	goarch: ["amd64", "arm64"]
	goos: ["darwin", "linux", "windows"]
	_flags: [
		"Version={{ .Tag }}",
		"Commit={{ .FullCommit }}",
		"BuildDate={{ .Date }}",
		"BuildOS={{ .Os }}",
		"BuildArch={{ .Arch }}",
		"BuildArm={{ .Arm }}",
	]
	ldflags: [ "-s -w", for f in _flags {"-X github.com/hofstadter-io/hof/cmd/hof/verinfo.\(f)"}]
	main: "main.go"
}]

dockers: [...{
	skip_push: false
}] & [
	// hof images
	for cfg in [
		{base: "debian", suf: ""},
		{base: "debian", suf: "debian-"},
		{base: "alpine", suf: "alpine-"},
	] {
		dockerfile: "../../ci/hof/docker/Dockerfile.\(cfg.base)"
		image_templates: [ for suf in ["{{.Tag}}", "{{ .ShortCommit }}", "latest"] {
			"ghcr.io/hofstadter-io/hof:\(cfg.suf)\(suf)"
		}]
	},
]

archives: [{
	// this makes it so a binary only is uploaded, rather than a tar file
	files: ["thisfiledoesnotexist*"]
	format:        "binary"
	name_template: "{{ .ProjectName }}_{{ .Tag }}_{{ .Os }}_{{ .Arch }}"
	//replacements: {
	//  amd64:   "x86_64"
	//  darwin:  "Darwin"
	//  linux:   "Linux"
	//  windows: "Windows"
	//}
}]

brews: [{
	name: "hof"

	homepage:    "https://github.com/hofstadter-io/hof"
	description: "CUE powered schemas, code gen, data modeling, dag engine, and tui."
	license:     "Apache-2"

	dependencies: [
		"docker",
	]

	download_strategy: "CurlDownloadStrategy"
	url_template:      "https://github.com/observIQ/observiq-otel-collector/releases/download/{{ .Tag }}/{{ .ArtifactName }}"
	folder:            "Formula"

	extra_install: #"""
		generate_completions_from_executable(bin / "hof", "completion")
		"""#

	commit_author: {
		name:  "dougbot"
		email: "bot@hofstadter.io"
	}

	commit_msg_template: "Brew formula update for {{ .ProjectName }} version {{ .Tag }}"

	skip_upload: false

	repository: {
		owner:  "hofstadter-io"
		name:   "homebrew-tap"
		branch: "master"
		token:  "{{ .Env.HOF_HOMEBREW_PAT }}"
		pull_request: {
			enabled: true
		}
	}
}]
