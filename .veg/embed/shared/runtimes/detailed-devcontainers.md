## Execution Environment

You run in an isolated environment using container technology.
Filesystem operation and command execution happen within this environment.
Your working directory is `basedir` from <env>.

<container>

https://github.com/qdm12/basedevcontainer

- `qmcgaw/basedevcontainer:debian` based on Debian Buster Slim in **376MB**
- All images are compatible with `amd64`, `386`, `arm64`, `armv7`, `armv6` and `ppc64le` CPU architectures
- Contains the packages:
  - `libstdc++`: needed by the VS code server
  - `zsh`: main shell instead of `/bin/sh`
  - `git`: interact with Git repositories
  - `openssh-client`: use SSH keys
  - `nano`: edit files from the terminal
- Contains the binaries:
  - [`gh`](https://github.com/cli/cli): interact with Github with the terminal
  - `docker`
  - `docker-compose` and `docker compose` docker plugin
  - [`docker buildx`](https://github.com/docker/buildx) docker plugin
  - [`bit`](https://github.com/chriswalz/bit)
  - [`devtainr`](https://github.com/qdm12/devtainr)
- Custom integrated terminal
  - Based on zsh and [oh-my-zsh](https://github.com/robbyrussell/oh-my-zsh)
  - Uses the [Powerlevel10k](https://github.com/romkatv/powerlevel10k) theme
  - With [Logo LS](https://github.com/Yash-Handa/logo-ls) as a replacement for `ls`
    - Shows information on login; easily extensible
- Cross platform
  - Easily bind mount your SSH keys to use with **git**
  - Manage your host Docker from within the dev container on Linux, MacOS and Windows
- Docker uses buildkit by default, with the latest Docker client binary.
- Extensible with docker-compose.yml
- Supports SSH keys with Linux, OSX and Windows

https://github.com/qdm12/godevcontainer

- `qmcgaw/godevcontainer:debian`
  - Based on Debian Buster Slim (size of 1.21GB)
- Based on [qmcgaw/basedevcontainer](https://github.com/qdm12/basedevcontainer)
  - Based on either Alpine or Debian
  - Minimal custom terminal and packages
  - See more [features](https://github.com/qdm12/basedevcontainer#features)
- Go 1.25 code obtained from the latest tagged Golang Docker image
- Go tooling [integrating with VS code](https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on), all cross built statically from source at the [binpot](https://github.com/qdm12/binpot):
  - [Google's Go language server gopls](https://github.com/golang/tools/tree/master/gopls)
  - [golangci-lint](https://github.com/golangci/golangci-lint), includes golint and other linters
  - [dlv](https://github.com/go-delve/delve/cmd/dlv) ⚠️ only works on `amd64` and `arm64`
  - [gomodifytags](https://github.com/fatih/gomodifytags)
  - [goplay](https://github.com/haya14busa/goplay)
  - [impl](https://github.com/josharian/impl)
  - [gotype-live](https://github.com/tylerb/gotype-live)
  - [gotests](https://github.com/cweill/gotests)
  - [gopkgs v2](https://github.com/uudashr/gopkgs/tree/master/v2)
- Terminal Go tools
  - [mockgen](https://github.com/golang/mock) to generate mocks
  - [mockery](https://github.com/vektra/mockery) to generate mocks for testify/mock
- Cross platform
  - Easily bind mount your SSH keys to use with **git**
  - Manage your host Docker from within the dev container, more details at [qmcgaw/basedevcontainer](https://github.com/qdm12/basedevcontainer#features)
- Extensible with docker-compose.yml
- Comes with extra Go binary tools for a few extra MBs: `kubectl`, `kubectx`, `kubens`, `stern` and `helm`

</container>
