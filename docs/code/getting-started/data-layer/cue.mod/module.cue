module: "hof.io/docs/example"
cue:    "0.5.0"

require: {
	"github.com/hofstadter-io/hof": "v0.6.8-beta.13"
}

indirect: {
	"github.com/hofstadter-io/ghacue":     "v0.2.0"
	"github.com/hofstadter-io/hofmod-cli": "v0.8.4"
}

replace: {
	"github.com/hofstadter-io/hof": "../../../../"
}
