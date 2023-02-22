package mod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/lib/mod/parse/modfile"
	"github.com/hofstadter-io/hof/lib/mod/parse/sumfile"
)

var upgradeMsg = `upgrading Hof CUE module
  cue.mods -> cue.mod/module.cue
  cue.sums -> cue.mod/sums.cue

You can now remove these files

`

func upgradeHofMods() error {
	_, err := os.Lstat("cue.mods")
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
			return err
		}
		// file does not exist
		return nil
	}

	fmt.Println(upgradeMsg)

	err = upgradeHofModFile()
	if err != nil {
		return err
	}

	err = upgradeHofSumFile()
	if err != nil {
		return err
	}

	// TODO, delete files after both have been updated & written

	return nil
}

func upgradeHofModFile() error {
	data, err := os.ReadFile("cue.mods")
	if err != nil {
		return err
	}

	mf, err := modfile.Parse("cue.mods", data, nil)
	if err != nil {
		return err
	}

	data, err = mf.WriteCUE()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("cue.mod/module.cue"), data, 0644)
}

func upgradeHofSumFile() error {
	data, err := os.ReadFile("cue.sums")
	if err != nil {
		return err
	}

	sf, err := sumfile.ParseSum(data, "cue.sums")
	if err != nil {
		return err
	}
	
	data, err = sf.WriteCUE()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("cue.mod/sums.cue"), data, 0644)
}
