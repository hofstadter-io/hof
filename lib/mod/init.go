package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gomod "golang.org/x/mod/module"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
)

var initFileContent = `module: %q
cue: %q
`

func Init(module string, rflags flags.RootPflagpole) (error) {
	upgradeHofMods()

	err := gomod.CheckPath(module)
	if err != nil {
		return fmt.Errorf("bad module name %q, should have domain format 'domain.com/...'", module)
	}

	_, err = os.Lstat("cue.mod")
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
			return err
		}
	} else {
		return fmt.Errorf("CUE module already exists in this directory")
	}

	s := fmt.Sprintf(initFileContent, module, verinfo.CueVersion)

	// mkdir & write file
	err = os.MkdirAll(filepath.Join("cue.mod", "pkg"), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join("cue.mod/module.cue"), []byte(s), 0644)
	if err != nil {
		return err
	}

	return nil
}

