package packer

import "time"

_vars: {
	_now:   string @tag(now,var=now)
	_out:   time.FormatString("20060102-150405", _now)
	suffix: string | *_out @tag(suffix)
}

_matrix: {
	tools: [
		"docker",
		"nerdctl",
		"nerdctl-rootless",
		"podman",
	]
	archs: [
		"amd",
		"arm",
	]
}

packer: [
	for tool in _matrix.tools for arch in _matrix.archs {
		"debian-\(tool)-\(arch).json"
	},
]

images: {
	for tool in _matrix.tools for arch in _matrix.archs {
		"\(tool)-\(arch)": _image & {#tool: tool, #arch: arch}
	}
}

_image: {
	#tool: string

	#arch: string

	// skip_create_image: true

	// foo
	builders: [{
		_name:        "debian-\(#tool)-\(#arch)"
		image_name:   "\(_name)-\(_vars.suffix)"
		image_family: _name
		type:         "googlecompute"
		zone:         "us-central1-a"
		project_id:   "hof-io--develop"

		ssh_username: "hof"

		disk_size: "25"
		disk_type: "pd-balanced"

		// base image
		_debianVersion: "v20230609"
		source_image:   [
				if #arch == "amd" {"debian-12-bookworm-\(_debianVersion)"},
				if #arch == "arm" {"debian-12-bookworm-arm64-\(_debianVersion)"},
				"unknown arch",
		][0]

		// machine type
		machine_type: [
				if #arch == "amd" {"n2-standard-2"},
				if #arch == "arm" {"t2a-standard-2"},
				"unknown arch",
		][0]
	}]
	provisioners: [
		{
			type: "shell"
			env: {
				"TOOL": #tool
				"ARCH": #arch
			}
			scripts: [
				"./scripts/packages.sh",
				"./scripts/\(#tool).sh",
			]
		},
	]

}
