package context

import (
	// "fmt"
	"path/filepath"

	"github.com/hofstadter-io/hof/pkg/ast"
	"github.com/hofstadter-io/hof/pkg/config"
)

func (ctx *Context) LoadModule(entrypoint string) (*ast.Module, error) {
	// TODO, look for module file, walk up if not found

	cfg, err := config.LoadModuleConfig(entrypoint)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("CFG: %#+v\n", cfg )

	mod := &ast.Module {
		Name: cfg.Name,
		Path: cfg.Path,
		Config: cfg,
		Packages: map[string]*ast.Package{},
	}
	ctx.Module = mod

	// first package to load
	dir := filepath.Join(entrypoint, cfg.Entrypoint)

	pkg, err := ctx.ReadPackage(dir, cfg)
	if err != nil {
		return mod, err
	}

	pkg.Path = cfg.Path

	mod.Packages[pkg.Path] = pkg

	return mod, nil
}

