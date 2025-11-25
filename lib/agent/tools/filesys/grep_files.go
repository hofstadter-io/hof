package filesys

import (
	"fmt"
	"os/exec"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type GrepFilesArgs struct {
	Path   string `json:"path"`   // the path to grep from
	Regexp string `json:"regexp"` // a regular expression to grep for
	// ExtraLines ?
}
type GrepFilesResult struct {
	Matches      string `json:"matches"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func NewGrepRegexp() (tool.Tool, error) {
	handler := func(ctx tool.Context, input GrepFilesArgs) (GrepFilesResult, error) {
		fmt.Println("grep_regexp.input:", input)
		var r GrepFilesResult

		// TODO
		// - parse regexp to make sure it is valid
		// - limit path, validate
		// - generally validation since we are effectively running on the user machine

		scriptFmt := `
		rg -Rn --sort=path -e '%s' %s
		`
		// -B%s -A%d  // for extra lines before / after
		// limit to a set of globs?

		script := strings.TrimSpace(fmt.Sprintf(scriptFmt, input.Regexp, input.Path))

		cmd := exec.Command("sh", "-c", script)
		if input.Path != "" {
			cmd.Dir = input.Path
		}
		out, err := cmd.CombinedOutput()
		if err != nil {
			r.Status = "error"
			r.ErrorMessage = fmt.Sprintf("Error:\n%s\n\nOutput:\n%s\n\n", err, string(out))
			return r, err
		}

		r.Status = "ok"
		r.Matches = string(out)
		return r, nil
	}
	return functiontool.New(functiontool.Config{
		Name:        "grep_regexp",
		Description: GrepRegexpDescription,
	}, handler)
}

const GrepRegexpDescription = `
Greps files for a regular expression from a path or the current directory if not set.
Matches are returned using this format:

path/to/file.ext
line_number:<content>'
line_number:<content>'
line_number:<content>'

path/file.ext
line_number:<content>'
line_number:<content>'
line_number:<content>'
`
