package visit

import (
	"errors"
	"fmt"

	"github.com/hofstadter-io/hof/pkg/ast"
)

func ToData(node ast.ASTNode) (interface{}, error) {

	switch N := node.(type) {
	case *ast.Module:
		return ModuleToData(N)

	case *ast.Package:
		return PackageToData(N)

	case *ast.Generator:
		return GeneratorToData(N)

	default:
		return nil, errors.New(fmt.Sprintf("Unknown ToData Type: %+v", node))
	}
}

func ModuleToData(N *ast.Module) (interface{}, error) {
	ret := map[string]interface{}{
		"type": "Module",
		"name": N.Name,
		"package": N.Path,
	}

	packages := map[string]interface{}{}
	ret["packages"] = packages

	for pkgName, pkg := range N.Packages {
		pkgData, err := ToData(pkg)
		if err != nil {
			return ret, err
		}
		packages[pkgName] = pkgData
	}

	return ret, nil
}

func PackageToData(N *ast.Package) (interface{}, error) {
	ret := map[string]interface{}{
		"type": "Package",
		"name": N.Name,
		"package": N.Path,
	}

	generators := map[string]interface{}{}
	ret["generators"] = generators

	for genName, gen := range N.Generators {
		genData, err := ToData(gen)
		if err != nil {
			return ret, err
		}
		generators[genName] = genData
	}

	return ret, nil
}

func GeneratorToData(N *ast.Generator) (interface{}, error) {
	ret := map[string]interface{}{
		"type": "Generator",
		"name": N.Name,
		"paths": N.Paths,
	}

	return ret, nil
}
