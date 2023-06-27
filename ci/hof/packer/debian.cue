package packer

import "time"

_vars: {
	_now:   string @tag(now,var=now)
	_out:   time.FormatString("20060102-150405", _now)
	suffix: string | *_out @tag(suffix)
}

_tools: [
	"docker",
	"nerdctl",
	"nerdctl-rootless",
	"podman",
]

images: {
	for tool in _tools {
		(tool): _image & {_tool: tool}
	}
}

_image: {
	_tool: string
	builders: [{
		image_name:   "debian-\(_tool)-\(_vars.suffix)"
		image_family: "hof-debian-\(_tool)"
		type:         "googlecompute"
		project_id:   "hof-io--develop"
		source_image: "debian-12-bookworm-v20230609"
		zone:         "us-central1-a"

		ssh_username: "hof"
		machine_type: "n2-standard-2"
		disk_size:    "25"
		disk_type:    "pd-ssd"

		// skip_create_image: true
	}]
	provisioners: [
		{
			type: "shell"
			scripts: [
				"./scripts/packages.sh",
				"./scripts/\(_tool).sh",
			]
		},
	]

}
