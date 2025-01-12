module: "github.com/hofstadter-io/hof"
source: kind: "git"
language: {
	version: "v0.9.0"
}
custom: {
	legacy: {
		cue: "0.10.0"
		require: {
			"github.com/hofstadter-io/cuelm":      "v0.1.1"
			"github.com/hofstadter-io/ghacue":     "v0.2.0"
			"github.com/hofstadter-io/hofmod-cli": "v0.9.0"
			"github.com/hofstadter-io/supacode":   "v0.0.7"
		}
	}
}
