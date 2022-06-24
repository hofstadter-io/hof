package gen

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

type RenderConfig struct {
	Data             cue.Value
	TemplateConfigs  []RenderTemplateConfig
	PartialFilenames []string
}

// parsed version of the --template flag
// semicolon separated: <filepath>;<?cuepath>;<?outpath>
type RenderTemplateConfig struct {
	// Filepath to the template
	Filepath string
	// 
	Cuepath  string
	Output   string 
}

type RenderEntry struct {
	Cuepath string
	TemplateFilename string

	Data cue.Value
}

func Render(args []string, cmdflags flags.RenderFlagpole) error {
	fmt.Println("Rendering:", args, cmdflags)

	// parse template flags

	// load CUE

	// process inputs

	return nil
}
