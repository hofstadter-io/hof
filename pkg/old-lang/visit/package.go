package visit

import (
	"errors"
	"fmt"

	"github.com/hofstadter-io/hof/pkg/old-lang/ast"
)

func CheckPackageFileDefines(pkg *ast.Package) error {
	// fmt.Println("CheckPackageFileDefines")

	packageName := ""
	packageFile := ""
	for fname, file := range pkg.Files {
		pname := file.PackageDecl.Name.Value

		if pname != packageName {
			if packageName == "" {
				packageName = pname
				packageFile = fname
			} else {
				msg := fmt.Sprintf(
					"Multiple packages defined in %s: (%s, %s) and (%s, %s)",
					pkg.Path,
					packageFile, packageName,
					fname, pname,
				)
				return errors.New(msg)
			}
		}

		/*
		fmt.Println(" -", fname)
		fmt.Printf("%#+v\n\n", file)
		*/
	}

	pkg.Name = packageName


	return nil
}
