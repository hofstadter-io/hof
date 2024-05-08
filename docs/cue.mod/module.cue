module: "github.com/hofstadter-io/hof-docs"
cue:    "v0.8.2"

require: {
	"github.com/hofstadter-io/cuelm": "v0.1.1"
	"github.com/hofstadter-io/hof":   "v0.6.9-rc.2"
}

indirect: {
	"github.com/hofstadter-io/ghacue":     "v0.2.0"
	"github.com/hofstadter-io/hofmod-cli": "v0.9.0"
	"github.com/hofstadter-io/supacode":   "v0.0.7"
}

replace: {
	"github.com/hofstadter-io/hof": "../"
}
