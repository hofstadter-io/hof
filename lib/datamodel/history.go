package datamodel

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/hofstadter-io/hof/lib/cuetils"
)

func HistoryBaseDir() (string, error) {
	dir, err := cuetils.FindModuleAbsPath("")
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, ".hof", "dm"), nil
}

func (V *Value) HistoryDir() (string, error) {
	basedir, err := HistoryBaseDir()
	if err != nil {
		return "", err
	}

	path := ""
	for v := V.Node; v != nil; v = v.Parent {
		id := v.Hof.Metadata.ID
		if id == "" {
			id = v.Hof.Metadata.Name
		}
		path = filepath.Join(id, path)
	}

	return filepath.Join(basedir, path), nil
}

// check for any history
func (dm *Datamodel) HasHistory() (bool, error) {
	return dm.T.HasHistory()
}

// check for any history
func (V *Value) HasHistory() (bool, error) {
	dir, err := V.HistoryDir()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(dir)
	if err != nil {
		// some other weir error
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return false, err
		}
		// case where there is no history
		return false, nil
	}

	return true, nil
}

/*
func (dm *Datamodel) LoadHistoryFrom(earliest string) error {
	// TODO, actually use earliest to flatten history at a timestamp or version

	return nil
}
*/

// LoadHistory loads the full datamodel history
func (dm *Datamodel) LoadHistory() error {

	has, err := dm.HasHistory()
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	err = dm.T.loadHistoryR()
	if err != nil {
		return err
	}

	return nil
}

func (V *Value) loadHistoryR() error {

	// load own history
	if V.Hof.Datamodel.History {
		err := V.loadHistory()
		if err != nil {

			return err
		}
	}

	// recurse if children to load any nested histories
	for _, c := range V.Children {
		err := c.T.loadHistoryR()
		if err != nil {
			return err
		}
	}

	return nil
}

func (V *Value) loadHistory() error {
	has, err := V.HasHistory()
	if err != nil || !has{
		return err
	}

	dir, err := V.HistoryDir()
	if err != nil {
		return err
	}

	fs, err := os.ReadDir(dir)
	if err != nil {
		return err
	}


	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		if ext != ".cue" {
			continue
		}
		s, err := loadSnapshot(dir, f.Name(), V.Value.Context())
		if err != nil {
			return err
		}

		s.Pos = len(V.history)
		V.history = append(V.history, s)
	}

	// sort by timestamp, most recent should be first
	sort.Slice(V.history, func(i,j int) bool {
		return V.history[i].Timestamp > V.history[j].Timestamp
	})
	

	return nil
}

// walks to see if there are any nodes with history as children
func (V *Value) hasHistBelow() bool {
	if V.Hof.Datamodel.History {
		return true
	}

	for _, c := range V.Children {
		b := c.T.hasHistBelow()
		if b {
			return b
		}
	}

	return false
}
