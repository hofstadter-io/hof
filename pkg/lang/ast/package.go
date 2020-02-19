package ast

type PackageDecl struct {
	// Parser filled
	BaseNode
	Name *Token
}

type Package struct {
	Parsed *PackageDecl
	BaseNode


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

	Generators map[string]*Generator
	Definitions map[string]*TypeDefinition

	PrivacyScopedNode
}

func NewPackage() *Package {
	return &Package {
		Files: map[string]*File{},
		Generators: map[string]*Generator{},
		Definitions: map[string]*TypeDefinition{},
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

