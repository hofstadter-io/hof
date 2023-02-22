module: "github.com/hof/test"
cue: "0.4.3"

require: {
	"github.com/hofstadter-io/hof": "v0.6.8-mod.1"
}

indirect: {
	"github.com/hofstadter-io/ghacue": "v0.2.0"
	"github.com/hofstadter-io/hofmod-cli": "v0.8.0"
}

replace: {
	"github.com/hofstadter-io/hof": "../../"
}
