package ast

type PackageDecl struct {
	// Parser filled
	ParseInfo *ParseInfo
	Name *Token
}

func (N *PackageDecl) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}


type Package struct {
	Name string
	// Full import string
	// github.com/hofstadter-io/hof-lang/modules/user
	Path string

	// If defined by a hof-lang.yaml

	// github.com
	Repo string
	// hofstadter-io
	Namespace string
	// hof-lang
	Package string
	// modules/user
	Subpath string

	Files map[string]*File

	SelfScope map[string]interface{}
	Exports   map[string]interface{}
}

func NewPackage() *Package {
	return &Package {
		Files: map[string]*File{},
		SelfScope: map[string]interface{}{},
		Exports: map[string]interface{}{},
	}
}

func (pkg *Package) AddFile(file *File) error {
	path := file.Path
	_, ok := pkg.Files[path]
	if ok {
		// already imported
	} else {
		pkg.Files[path] = file
	}
	return nil

}

