package context

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/pkg/lang/ast"
	"github.com/hofstadter-io/hof/pkg/lang/module"
)


func (ctx *Context) ReadPackage(dir string, cfg *module.Config) (*ast.Package, error) {
	// fmt.Println("\nReadPackage:", dir, cfg.Entrypoint, cfg.Path)

	var rootModule = false
	if dir == cfg.Entrypoint {
		// Do nothing because this is the first package to be read
	} else if strings.HasPrefix(dir, cfg.Path) {
		// Check for within module import
		rootModule = true
		subdir := strings.TrimPrefix(dir, cfg.Path)
		dir = filepath.Join(cfg.Entrypoint, subdir)
	} else {
		// Assume this is a vendor package, add prefix: <entrypoint>/vendor/
		dir = filepath.Join(cfg.Entrypoint, "vendor", dir)
	}

	// fmt.Println("DynamicDir:", dir)
	epkg, ok := ctx.Packages[dir]
	if ok {
		return epkg, nil
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	pkg := ast.NewPackage()
	pkg.Path = dir
	// Are we just undoing here? (does not seem like it, because we use dir below for filesystem path, and this is just a "import path")
	if !strings.HasPrefix(pkg.Path, "vendor") {
		pkg.Path = filepath.Join(ctx.Module.Path, pkg.Path)
		// fmt.Printf("PKG: %#+v\n", *file.Package)
	}

	ctx.AddPackage(pkg)
	if rootModule {
		ctx.Module.AddPackage(pkg)
	}

	for _, info := range infos {
		// only want hof files from this directory
		if info.IsDir() {
			continue
		}
		if filepath.Ext(info.Name()) != ".hof" {
			continue
		}

		// read the hof file
		hofFilePath := filepath.Join(dir, info.Name())
		// fmt.Println(" ~", hofFilePath)

		file, err := ctx.ReadFile(hofFilePath, cfg)
		if err != nil {
			ctx.AddError(err)
			pkg.AddFile(&ast.File{
				Name: info.Name(),
				Path: hofFilePath,
			})
			continue
		}

		file.Name = info.Name()
		file.Package = pkg
		pkg.AddFile(file)

		err = ctx.LoadFileImports(file)
		if err != nil {
			ctx.AddError(err)
		}

	}

	return pkg, nil
}

