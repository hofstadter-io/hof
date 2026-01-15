module: "github.com/hofstadter-io/hof"
cue: "0.5.0"

// Direct dependencies (managed by hof)
require: {
	"github.com/hofstadter-io/ghacue": "v0.2.0"
	"github.com/hofstadter-io/hofmod-cli": "v0.8.4"
}

// Indirect dependencies (managed by hof)
indirect: { ... }

// Replace dependencies with remote or local, useful for co-development
replace: {
  "github.com/hofstadter-io/ghacue": "github.com/myorg/ghacue": "v0.2.1"
  "github.com/hofstadter-io/hofmod-cli": "../mods/cli"
}
