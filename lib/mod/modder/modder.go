package modder

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/yagu"
)

// This modder is for more complex, yet configurable module processing.
// You can have system wide and local custom configurations.
// The fields in this struct are alpha and are likely to change
type Modder struct {
	// MetaConfiguration
	Name    string `yaml:"Name"`
	Version string `yaml:"Version",omitempty`

	// Module information
	ModFile       string `yaml:"ModFile",omitempty`
	SumFile       string `yaml:"SumFile",omitempty`
	ModsDir       string `yaml:"ModsDir",omitempty`
	MappingFile   string `yaml:"MappingFile",omitempty`
	PrivateEnvVar string `yaml:"PrivateEnvVar",omitempty`

	// Commands override default, configuragble processing
	// for things like golang
	NoLoad        bool       `yaml:"NoLoad",omitempty`
	CommandInit   [][]string `yaml:"CommandInit",omitempty`
	CommandGraph  [][]string `yaml:"CommandGraph",omitempty`
	CommandTidy   [][]string `yaml:"CommandTidy",omitempty`
	CommandVendor [][]string `yaml:"CommandVendor",omitempty`
	CommandVerify [][]string `yaml:"CommandVerify",omitempty`
	CommandStatus [][]string `yaml:"CommandStatus",omitempty`

	// Init related fields
	// we need to create things like directories and files beyond the
	InitTemplates    map[string]string `yaml:"InitTemplates",omitempty`
	InitPreCommands  [][]string        `yaml:"InitPreCommands",omitempty`
	InitPostCommands [][]string        `yaml:"InitPostCommands",omitempty`

	// Vendor related fields
	// filesystem globs for discovering files we should copy over
	VendorIncludeGlobs []string `yaml:"VendorIncludeGlobs",omitempty`
	VendorExcludeGlobs []string `yaml:"VendorExcludeGlobs",omitempty`
	// Any files we need to generate
	VendorTemplates    map[string]string `yaml:"VendorTemplates",omitempty`
	VendorPreCommands  [][]string        `yaml:"VendorPreCommands",omitempty`
	VendorPostCommands [][]string        `yaml:"VendorPostCommands",omitempty`

	// Some more vendor controls
	ManageFileOnly       bool `yaml:"ManageFileOnly",omitempty`
	SymlinkLocalReplaces bool `yaml:"SymlinkLocalReplaces",omitempty`

	// Introspection Configuration(s)
	// filesystem globs for discovering files we should introspect
	// regexs for extracting package information
	IntrospectIncludeGlobs []string `yaml:"IntrospectIncludeGlobs",omitempty`
	IntrospectExcludeGlobs []string `yaml:"IntrospectExcludeGlobs",omitempty`
	IntrospectExtractRegex []string `yaml:"IntrospectExtractRegex",omitempty`

	PackageManagerDefaultPrefix string `yaml:"PackageManagerDefaultPrefix",omitempty`

	// filesystem
	FS billy.Filesystem `yaml:"-"`

	// root module
	workdir string	`yaml:"-"`
	module  *Module `yaml:"-"`
	errors  []error `yaml:"-"`

	// dependency modules (requires/replace)
	// dependencies should respect any .mvsconfig it finds along side the module files
	// module writers can then have local control over how their module is handeled during vendoring
	depsMap map[string]*Module `yaml:"-"`

	// compiled cue, used for merging
	CueInstance *cue.Instance `yaml:"-"`
}

func NewFromFile(lang, filepath string, FS billy.Filesystem) (*Modder, error) {

	bytes, err := yagu.BillyReadAll(filepath, FS)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return nil, err
		}
		// The user has not setup a global $HOME/.mvs/mvsconfig file
		return nil, nil
	}

	var mdrMap map[string]*Modder

	rt := cue.Runtime{}
	i, err := rt.Compile(filepath, string(bytes))
	if err != nil {
		return nil, err
	}
	err = i.Value().Decode(&mdrMap)
	if err != nil {
		return nil, err
	}

	mdr, ok := mdrMap[lang]
	if !ok {
		return nil, fmt.Errorf("lang %q not found in %s", lang, filepath)
	}

	return mdr, nil
}
