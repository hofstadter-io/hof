package datamodel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "github.com/Masterminds/semver/v3"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
)

// YYYYMMDDHHMMSS in Golang
const CheckpointFmt = "20060102150405"

// Snapshot represents a Value at a point in time
type Snapshot struct {
	// The current value at this snapshot
	Value cue.Value

	// Position in history
	Pos int

	// Point it was snapshotted
	Timestamp string

	// crypto hash of (snapshot data?)
	Hash string

	// Explination for the snapshot or changes therein
	Message string

	// Lenses between neighboring snapshots
	Lense Lense
}

func (V *Value) LatestSnapshot() *Snapshot {
	if len(V.history) == 0 {
		return nil
	}

	return V.history[0]
}

func loadSnapshot(dir, fpath string, ctx *cue.Context) (*Snapshot, error) {
	_, fname := filepath.Split(fpath)
	ts := strings.TrimPrefix(fname, "ver_")
	ts = strings.TrimSuffix(ts, ".cue")
	s := &Snapshot{
		Timestamp: ts,
	}

	fn := filepath.Join(dir,fpath)
	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	v := ctx.CompileBytes(data, cue.Filename(fn))
	if v.Err() != nil {
		return nil, v.Err()
	}
	s.Value = v

	msg := fmt.Sprintf("msg_%s", ts)
	cmsg := v.LookupPath(cue.ParsePath(msg))
	s.Message, err = cmsg.String()
	if err != nil {
		return nil, err
	}

	return s, nil
}

// we probably want this to be closer to `cue def`
// with imports, schemas, and such... indep eval'able
func writeSnapshot(dir, fname, message, pkgId string, V cue.Value) error {
	ver := fmt.Sprintf("ver_%s", fname)
	msg := fmt.Sprintf("msg_%s", fname)

	err := V.Validate(cue.Concrete(true))
	if err != nil {
		return err
	}

	// build a new top-level cue label and value
	val := V.Context().CompileString("_")
	val = val.FillPath(cue.ParsePath(msg), message)
	val = val.FillPath(cue.ParsePath(ver), V)

	node := val.Syntax(
		cue.Final(),
		cue.Docs(true),
		cue.Attributes(true),
		cue.Definitions(true),
		cue.Optional(true),
		cue.Hidden(true),
		cue.Concrete(true),
		cue.ResolveReferences(true),
	)

	file, err := astutil.ToFile(node.(*ast.StructLit))
	if err != nil {
		return err
	}

	pkg := &ast.Package{
		Name: ast.NewIdent(pkgId),
	}
	file.Decls = append([]ast.Decl{pkg}, file.Decls...)

	// fmt.Printf("%#+v\n", file)

	bytes, err := format.Node(
		file,
		format.Simplify(),
	)
	if err != nil {
		return err
	}

	// make history dir
	err = yagu.Mkdir(dir)
	if err != nil {
		return err
	}

	// write file
	str := string(bytes)
	fn := filepath.Join(dir, fmt.Sprintf("%s.cue", fname))
	err = os.WriteFile(fn, []byte(str), 0666)
	if err != nil {
		return err
	}

	return nil
}

// MakeSnapshot creates a new snapshot for each history annotation in the Datamodel tree
func (dm *Datamodel) MakeSnapshot(timestamp string, dflags flags.DatamodelPflagpole, cflags flags.Datamodel__CheckpointFlagpole) error {
	err := dm.T.makeSnapshotR(timestamp, cflags.Message)
	if err != nil {
		return err
	}

	return nil
}

func (V *Value) makeSnapshotR(timestamp, message string) error {

	// load own history
	if V.Hof.Datamodel.History {
		err := V.makeSnapshot(timestamp, message)
		if err != nil {
			return err
		}
	}

	// recurse if children to load any nested histories
	for _, c := range V.Children {
		err := c.T.makeSnapshotR(timestamp, message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (V *Value) makeSnapshot(timestamp, message string) error {

	// if no diff, no snapshot
	if len(V.history) > 0 && !V.hasDiff() {
		return nil
	}

	dir, err := V.HistoryDir()
	if err != nil {
		return err
	}

	val := V.Value

	// enrich val
	// - $hof
	// - lacunas

	pkg := V.Hof.Metadata.ID
	if pkg == "" {
		pkg = V.Hof.Metadata.Name
	}
	pkg = strings.Replace(pkg, "-", "_", -1)
	err = writeSnapshot(dir, timestamp, message, pkg, val)
	if err != nil {
		return err
	}

	return nil
}



/*  HAS OLD BUMP LOGIC
func checkpointDatamodel(dm *Datamodel, timestamp, bump string) error {
	// check subsumption
	if dm.Subsume != nil {
		fmt.Println("subsume:", dm.Subsume)
		if bump == "patch" {
			return fmt.Errorf("\nBackwards incompatible changes, you should use a major or minor bump")
		}
	}

	tag := "0.0.0"

	if len(dm.History.Past) > 0 {
		// get last semver tag
		tag = dm.History.Past[0].Version
	}

	v, err := semver.StrictNewVersion(tag)
	if err != nil {
		return fmt.Errorf("error parsing last semver in history: %q\n%s", tag, err.Error())
	}

	// bump semver tag
	switch bump {
	case "major":
		b := v.IncMajor()
		v = &b
	case "minor":
		b := v.IncMinor()
		v = &b
	case "patch":
		b := v.IncPatch()
		v = &b
	default:
		n, err := semver.StrictNewVersion(bump)
		if err != nil {
			return fmt.Errorf("error parsing next semver: %q\n%s", bump, err.Error())
		}
		if !n.GreaterThan(v) {
			return fmt.Errorf("new version must be greater than last: %v < %v", n, v)
		}
		v = n
	}

	tag = v.String()
	dm.Version = tag

	fmt.Printf("%s bump ver to: %s\n", bump, dm.Version)

	// add to dm
	ctx := dm.Value.Context()
	str := fmt.Sprintf("{ @dm_ver(%s) }", tag)
	val := ctx.CompileString(str)
	dm.Value = dm.Value.Unify(val)

	return writeCheckpoint(dm, timestamp)
}
*/

// GetSnapshotList finds all snapshots from top-level @history() values.
// It returns a map where the key is the path to the @history
// and the value is a list of all snapshots at that value.
// The top most parent will contain all snapshots for any nested values,
// due to the way diffs propagate up the hof node tree
// You should not modify the returned snapshots
func (dm *Datamodel) GetSnapshotList() (map[string][]*Snapshot) {
	return dm.T.getSnapshotListR()
}

func (V *Value) getSnapshotListR() (map[string][]*Snapshot) {
	// if we find a history point, it should be the top-level
	// so stop recursion and return the snapshot list at this node
	if V.Hof.Datamodel.History {
		return V.getSnapshotList()
	}

	// otherwise recurse over children to build up a larger map
	ret := make(map[string][]*Snapshot)
	for _, c := range V.Children {
		R := c.T.getSnapshotListR()
		for k,v := range R {
			ret[k] = v
		}
	}

	return ret
}

func (V *Value) getSnapshotList() (map[string][]*Snapshot) {
	return map[string][]*Snapshot {
		V.Hof.Path: V.history,
	}
}

func (V *Value) hasSnapshotAt(timestamp string) bool {
	for _, s := range V.history {
		if s.Timestamp == timestamp {
			return true
		}
	}
	return false
}

func (V *Value) getSnapshotPos(timestamp string) int {
	for i, s := range V.history {
		if s.Timestamp == timestamp {
			return i
		}
	}
	return -1
}
