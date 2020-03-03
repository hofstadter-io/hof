package visit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/pkg/lang/ast"
)

func ToData(node ast.ASTNode) (interface{}, error) {

	switch N := node.(type) {
	case *ast.Module:
		return ModuleToData(N)

	case *ast.Package:
		return PackageToData(N)

	case *ast.TypeDefinition:
		return DefinitionToData(N)

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

	definitions := map[string]interface{}{}
	ret["definitions"] = definitions

	for defName, def := range N.Definitions {
		defData, err := ToData(def)
		if err != nil {
			return ret, err
		}
		definitions[defName] = defData
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

func DefinitionToData(N *ast.TypeDefinition) (interface{}, error) {
	paths := []string{}
	for _, p := range N.Paths {
		ps := strings.Join(p.Paths, ".")
		paths = append(paths, ps)
	}

	ret := map[string]interface{}{
		"type": "TypeDefintion",
		"name": N.Name.Value,
		"paths": paths,
	}

	return ret, nil
}

func GeneratorToData(N *ast.Generator) (interface{}, error) {
	paths := []string{}
	for _, p := range N.Paths {
		ps := strings.Join(p.Paths, ".")
		paths = append(paths, ps)
	}

	ret := map[string]interface{}{
		"type": "Generator",
		"name": N.Name,
		"paths": paths,
	}

	return ret, nil
}
