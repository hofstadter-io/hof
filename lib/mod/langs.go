package mod

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	"github.com/hofstadter-io/hof/lib/mod/langs"
	"github.com/hofstadter-io/hof/lib/mod/modder"
)

// FROM the USER's HOME dir
const GLOBAL_MVS_CONFIG = "hof/.mvsconfig.cue"
const LOCAL_MVS_CONFIG = ".mvsconfig.cue"

var (
	// Default known modderr
	LangModderMap = langs.DefaultModders
)

const knownLangMessage = `
Known Languages:

  %s

For more info on a language:

  hof mod info <lang>
`

func DiscoverLangs() (langs []string) {

	for lang, mdr := range LangModderMap {
		// Let's check for a custom
		_, err := os.Lstat(mdr.ModFile)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
				fmt.Println(err)
				// return err // more of a warning right now
			}
			// file not found
			continue
		}
		// we found a mod file
		langs = append(langs, lang)
	}

	return langs
}

func KnownLangs() string {
	// extract and sort for consistency
	langs := []string{}

	for lang, _ := range LangModderMap {
		langs = append(langs, lang)
	}

	sort.Strings(langs)
	langStr := strings.Join(langs, "\n  ")

	msg := fmt.Sprintf(knownLangMessage, langStr)

	return msg
}

const unknownLangMessage = `
Unknown language %q.

Please check the following files for definitions
  %s  (in the current directory)
	$HOME/%s

To see a list of known languages from the current directory:

  hof mod info
`

func LangInfo(lang string) (string, error) {

	if lang == "" {
		return KnownLangs(), nil
	}

	modder, ok := LangModderMap[lang]
	if !ok {
		return "", fmt.Errorf(unknownLangMessage, lang, LOCAL_MVS_CONFIG, GLOBAL_MVS_CONFIG)
	}

	// fmt.Printf("=====\n%#+v\n=====\n", modder)

	// TODO output as cue
	bytes, err := yaml.Marshal(modder)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// InitLangs loads the language modders in the following order
//   1. Defaults compiled into hof
//   2. Globals found in <user-config>/hof/.mvsconfig.cue
//   3. Local (CWD)/.mvsconfig
//
// New languages discovered are appended, existing are replaced completely
func InitLangs() {
	var err error

	// combine the specs with the default lang modder defs
	defaultStr := langs.ModderSpec + langs.DefaultLangs

	// Read and parse in the default langs specs
	ctx := cuecontext.New()
	dLangs := ctx.CompileString(defaultStr, cue.Filename("langs.cue"))
	err = dLangs.Err()
	if err != nil {
		panic(err)
	}

	// Validate the langs
	err = dLangs.Value().Validate()
	if err != nil {
		panic(err)
	}

	// decode defaults into a temp map
	var mdrMap map[string]*modder.Modder
	err = dLangs.Value().LookupPath(cue.ParsePath("langs")).Decode(&mdrMap)
	if err != nil {
		panic(err)
	}

	// load from temp into lang.{Default,Loaded}Modders
	for lang, mdr := range mdrMap {
		langs.DefaultModders[lang] = mdr
		langs.LoadedModders[lang] = mdr
		LangModderMap[lang] = mdr
	}

	// Global Language Modder Config
	homedir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
	}
	err = initFromFile(path.Join(homedir, GLOBAL_MVS_CONFIG))
	if err != nil {
		fmt.Println(err)
	}

	// Local Language Modder Config
	err = initFromFile(LOCAL_MVS_CONFIG)
	if err != nil {
		fmt.Println(err)
	}
}

func initFromFile(filepath string) error {
	// Reand an MVS config file (cue format)
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return err
		}
		// The user has not setup a global $HOME/.mvs/mvsconfig file
		return nil
	}

	mdrStr := langs.ModderSpec + string(bytes)

	var mdrMap map[string]*modder.Modder

	// Compile the config into cue
	ctx := cuecontext.New()
	i := ctx.CompileString(mdrStr, cue.Filename(filepath))
	err = i.Err()
	if err != nil {
		return err
	}
	err = i.Value().LookupPath(cue.ParsePath("langs")).Decode(&mdrMap)
	if err != nil {
		return err
	}

	for lang, mdr := range mdrMap {
		langs.LoadedModders[lang] = mdr
		LangModderMap[lang] = mdr
	}

	return nil
}
