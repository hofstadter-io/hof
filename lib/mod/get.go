package mod

import (
	"fmt"
	"strings"

	gomod "golang.org/x/mod/module"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/repos/cache"
)

func Get(module string, rflags flags.RootPflagpole, mflags flags.ModPflagpole) (error) {
	upgradeHofMods()

	parts := strings.Split(module, "@")
	if len(parts) != 2 {
		return fmt.Errorf("bad module format %q, should be 'domain.com/path/...@version'", module)
	}

	path, ver := parts[0], parts[1]

	err := gomod.CheckPath(path)
	if err != nil {
		return fmt.Errorf("bad module name %q, should have domain format 'domain.com/...'", path)
	}

	_, err = cache.Cache(path, ver)
	if err != nil {
		return err
	}

	cm, err := loadRootMod()
	if err != nil {
		return err
	}

	if path == cm.Module {
		return fmt.Errorf("cannot get current module")
	}

	// check for indirect and delete
	_, ok := cm.Indirect[path]
	if ok {
		delete(cm.Indirect, path)
	}

	_, ok = cm.Require[path]
	if ok {
		return fmt.Errorf("module already required")
	}

	// add dep to required
	cm.Require[path] = ver

	fns := []func () error {
		cm.SolveMVS,
		cm.CleanDeps,
		cm.CleanSums,
		cm.WriteModule,
	}

	for _, fn := range fns {
		err := fn()
		if err != nil {
			return err
		}
	}

	return cm.Vendor(true, rflags.Verbosity)
}

