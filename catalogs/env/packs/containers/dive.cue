package containers

import "github.com/hofstadter-io/hof/schemas/env"

dive: {
	#ver: string | *"0.13.1"

	#arch:   string | *"arm64" | "amd64" // todo, this should default to current OS
	#distro: string | *"linux"

	cli: {
		// hack for now to maintain consistency in pack.<tool>.cli.install
		install: env.Sh & {
			_file:  "dive_\(#ver)_\(#distro)_\(#arch).tar.gz"
			_ibin:  "/usr/local/bin/cosign"
			_src:   "https://github.com/wagoodman/dive/releases/download/v\(#ver)/\(_file)"
			script: """
        cd /tmp
        wget -q \(_src)
        tar -xzf \(_file)
        ls -lh *
        mv dive /usr/local/bin/
        rm -rf /tmp/*
        """
		}

	}
}
