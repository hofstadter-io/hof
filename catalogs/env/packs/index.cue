package packs

import (
	icontainers "github.com/hofstadter-io/hof/catalogs/env/packs/containers"
	idatabases "github.com/hofstadter-io/hof/catalogs/env/packs/databases"
	ilang "github.com/hofstadter-io/hof/catalogs/env/packs/lang"
	itool "github.com/hofstadter-io/hof/catalogs/env/packs/tool"
)

flags: {
	goos: string
	arch: string
}

containers: icontainers & { "flags": flags }
databases:  idatabases & { "flags": flags }
lang:       ilang & { "flags": flags }
tool:       itool & { "flags": flags }
