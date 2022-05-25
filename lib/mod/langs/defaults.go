package langs

import (
	"github.com/hofstadter-io/hof/lib/mod/modder"
)

var DefaultModders = make(map[string]*modder.Modder)
var LoadedModders = make(map[string]*modder.Modder)

var DefaultModdersCue = map[string]string{
	"go":     GolangModder,
	"cue":    CuelangModder,
	"python": PythonModder,
}

var DefaultLangs = CuelangModder +
	GolangModder +
	PythonModder
