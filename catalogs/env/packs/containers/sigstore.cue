package containers

import "github.com/hofstadter-io/hof/schemas/env"

cosign: {
	#ver: string | *"3.0.3"

	#arch:   string | *"arm64" | "amd64" // todo, this should default to current OS
	#distro: string | *"linux"

	cli: {
		// hack for now to maintain consistency in pack.<tool>.cli.install
		install: env.Sh & {
			_file:  "cosign-\(#distro)-\(#arch)"
			_ibin:  "/usr/local/bin/cosign"
			_src:   "https://github.com/sigstore/cosign/releases/download/v\(#ver)/\(_file)"
			script: "wget -q \(_src) -O \(_ibin) && chmod +x \(_ibin)"
		}

	}
}
