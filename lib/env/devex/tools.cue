package devex

import (
	"strings"

	"github.com/hofstadter-io/hof/schemas/env"
)


// hof/veg

_tools: {
	// todo, eventually we can use this dag stuff to redo how cue, go, and others work
	// indep container -> file/dir -> building container
	// cue: env.File & {
	// 	path: "/cue"
	// 	source: env.Container & {}
	// }

	githubBin: env.Exec & {
		#ver:  string
		#arch: *"arm64" | "amd64" | string
    #distro: string | *"linux"

		#repo: string
		#name: string | *strings.Split(#repo, "/")[1]
		#bins: [...string] | *[#name]
		_bins: strings.Join(#bins, " ")

		args: ["sh", "-c", _script]

		_file:   "\(#name)_v\(#ver)_\(#distro)_\(#arch).tar.gz"
		_src:    "https://github.com/\(#repo)/releases/download/v\(#ver)/\(_file)"
		_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    tar -xzf \(_file)
    mv \(_bins) /usr/local/bin/
    rm -rf /tmp/*
    """
	}

	hashicorpBin: env.Exec & {
		#ver:  string | *"1.14.3"
		#arch: *"arm64" | "amd64"

		#tool: string

		args: ["sh", "-c", _script]
		_file:   "\(#tool)_\(#ver)_linux_\(#arch).zip"
		_src:    "https://releases.hashicorp.com/\(#tool)/\(#ver)/\(_file)"
		_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    unzip \(_file)
    mv \(#tool) /usr/local/bin/\(#tool)
    rm -rf /tmp/*
    """
	}

	go: [
    // set Go ENV vars before installing, especially extra packages
    env.Env & {PATH: "$PATH:/usr/local/go/bin"},
    env.Env & {GOBIN: "/usr/local/bin"}, // go install to /usr/local/bin
    env.Env & {GOPATH: "/go"}, // todo, dagger cache
    env.Env & {GOCACHE: "/cache/go"}, // todo, dagger cache
    env.Exec & {
      #ver:  string | *"1.25.4"
      #arch: *"arm64" | "amd64"
      args: ["sh", "-c", _script]

      _file:   "go\(#ver).linux-\(#arch).tar.gz"
      _src:    "https://go.dev/dl/\(_file)"
      _script: """
      set -eou pipefail

      cd /tmp
      wget -q \(_src)
      tar -C /usr/local -xzf \(_file)
      rm -rf /tmp/*

      # LSP
      go install golang.org/x/tools/gopls@latest

      # lint tools
      go install honnef.co/go/tools/cmd/staticcheck@latest
      go install github.com/mgechev/revive@latest
      curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b /usr/local/bin v2.7.2
      """
    },
  ]

  node: env.Exec & {
    #ver: string | *"24.12.0"
		#arch: *"arm64" | "amd64"

    // # NOT a "standalone" binary as the say (it requires includes and more, so by def not standalone... fucktards...)
    _src: "https://nodejs.org/dist/v\(#ver)/\(_file)"
    _file: "node-v\(#ver)-linux-\(#arch).tar.xz"

		args: ["sh", "-c", _script]

		_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    tar -C /usr/local -xf \(_file) --strip-components=1

    # package manager
    corepack enable pnpm

    # LSP
    npm install -g tsx typescript typescript-language-server
    """
  }

  python: env.Exec & {
		args: ["sh", "-c", _script]

    // yes, pyright requires node and recommends installing via npm
		_script: """
    set -eou pipefail

    # uv 
    curl -LsSf https://astral.sh/uv/install.sh | env UV_INSTALL_DIR="/usr/local/bin" sh

    # LSP
    npm install -g pyright
    """
  }

  agents: {
    lsp2mcp: env.Exec & {
      args: ["sh", "-c", _script]

      _script: """
      # LSP -> MCP
      go install github.com/isaacphi/mcp-language-server@latest
      """
    }
  }

	zsh: env.Exec & {
		args: ["sh", "-c", _script]
		_script: """
			sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
			sed -i 's/^ZSH_THEME=.*/ZSH_THEME="frisk"/' /root/.zshrc
			"""
	}

	dockerRepo: env.Exec & {
		args: ["sh", "-c", _script]

		_script: """
			set -eou pipefail

			install -m 0755 -d /etc/apt/keyrings
			curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
			chmod a+r /etc/apt/keyrings/docker.asc


			cat <<EOF > /etc/apt/sources.list.d/docker.sources
			Types: deb
			URIs: https://download.docker.com/linux/debian
			Suites: trixie
			Components: stable
			Signed-By: /etc/apt/keyrings/docker.asc
			EOF
			"""
	}

  kubectl: env.Exec & {
    #ver: string | *"1.31.0"
		#arch: *"arm64" | "amd64"

		args: ["sh", "-c", _script]

    _src: "https://storage.googleapis.com/kubernetes-release/release/v\(#ver)/bin/linux/\(#arch)/kubectl"

		_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    chmod +x ./kubectl
    mv kubectl /usr/local/bin/kubectl
    rm -rf /tmp/*
    """

  }

  helm: env.Exec & {
    #ver: string | *"4.0.4"
		#arch: *"arm64" | "amd64"

		args: ["sh", "-c", _script]

    _src: "https://get.helm.sh/\(_file)"
    _file: "helm-v\(#ver)-linux-\(#arch).tar.gz"

		_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    tar -xzf \(_file)
    mv linux-\(#arch)/helm /usr/local/bin/helm
    rm -rf /tmp/*
    """
  }

  // basically the same as github, but without the version in teh filname
	crane: env.Exec & {
		#ver:  string | *"0.20.7"
		#arch: string | *"arm64" | "x86_64"
    #distro: string | *"Linux"

    #repo: "google/go-containerregistry"
    #bins: ["crane", "krane"]

		#name: string | *strings.Split(#repo, "/")[1]
		_bins: strings.Join(#bins, " ")

		args: ["sh", "-c", _script]

		_file:   "\(#name)_\(#distro)_\(#arch).tar.gz"
		_src:    "https://github.com/\(#repo)/releases/download/v\(#ver)/\(_file)"
		_script: """
    set -eou pipefail

    cd /tmp
    wget -q \(_src)
    tar -xzf \(_file)
    mv \(_bins) /usr/local/bin/
    rm -rf /tmp/*
    """
	}

}
