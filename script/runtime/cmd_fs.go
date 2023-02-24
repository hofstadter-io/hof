package runtime

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// exists
// symlink
// cp
// mkdir
// rm
// chmod
// chown

// cd changes to a different directory.
func (ts *Script) CmdCd(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? cd")
	}
	if len(args) != 1 {
		ts.Fatalf("usage: cd dir")
	}

	dir := args[0]
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(ts.cd, dir)
	}
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		ts.Fatalf("directory %s does not exist", dir)
	}
	ts.Check(err)
	if !info.IsDir() {
		ts.Fatalf("%s is not a directory", dir)
	}
	ts.cd = dir
	ts.Logf("%s\n", ts.cd)
}

func (ts *Script) CmdChmod(neg int, args []string) {
	if len(args) != 2 {
		ts.Fatalf("usage: chmod mode file")
	}
	mode, err := strconv.ParseInt(args[0], 8, 32)
	if err != nil {
		ts.Fatalf("bad file mode %q: %v", args[0], err)
	}
	if mode > 0777 {
		ts.Fatalf("unsupported file mode %.3o", mode)
	}
	err = os.Chmod(ts.MkAbs(args[1]), os.FileMode(mode))
	if neg > 0 {
		if err == nil {
			ts.Fatalf("unexpected chmod success")
		}
		return
	}
	if err != nil {
		ts.Fatalf("unexpected chmod failure: %v", err)
	}
}

// cp copies files, maybe eventually directories.
func (ts *Script) CmdCp(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? cp")
	}
	if len(args) < 2 {
		ts.Fatalf("usage: cp src... dst")
	}

	dst := ts.MkAbs(args[len(args)-1])
	info, err := os.Stat(dst)
	dstDir := err == nil && info.IsDir()
	if len(args) > 2 && !dstDir {
		ts.Fatalf("cp: destination %s is not a directory", dst)
	}

	for _, arg := range args[:len(args)-1] {
		var (
			src  string
			data []byte
			mode os.FileMode
		)
		switch arg {
		case "stdout", "@stdout":
			src = arg
			data = []byte(ts.stdout)
			mode = 0666
		case "stderr", "@stderr":
			src = arg
			data = []byte(ts.stderr)
			mode = 0666
		default:
			src = ts.MkAbs(arg)
			info, err := os.Stat(src)
			ts.Check(err)
			mode = info.Mode() & 0777
			data, err = ioutil.ReadFile(src)
			ts.Check(err)
		}
		targ := dst
		if dstDir {
			targ = filepath.Join(dst, filepath.Base(src))
		}
		ts.Check(ioutil.WriteFile(targ, data, mode))
	}
}

// exists checks that the list of files exists.
func (ts *Script) CmdExists(neg int, args []string) {
	var readonly bool
	if len(args) > 0 && args[0] == "-readonly" {
		readonly = true
		args = args[1:]
	}
	if len(args) == 0 {
		ts.Fatalf("usage: exists [-readonly] file...")
	}

	for _, file := range args {
		file = ts.MkAbs(file)
		info, err := os.Lstat(file)
		if err == nil && neg > 0 {
			what := "file"
			if info.IsDir() {
				what = "directory"
			}
			ts.Fatalf("%s %s unexpectedly exists", what, file)
		}
		if err != nil && neg == 0 {
			ts.Fatalf("%s does not exist", file)
		}
		if err == nil && neg == 0 && readonly && info.Mode()&0222 != 0 {
			ts.Fatalf("%s exists but is writable", file)
		}
	}
}

// mkdir creates directories.
func (ts *Script) CmdMkdir(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? mkdir")
	}
	if len(args) < 1 {
		ts.Fatalf("usage: mkdir dir...")
	}
	for _, arg := range args {
		ts.Check(os.MkdirAll(ts.MkAbs(arg), 0777))
	}
}

// rm removes files or directories.
func (ts *Script) CmdRm(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? rm")
	}
	if len(args) < 1 {
		ts.Fatalf("usage: rm file...")
	}
	for _, arg := range args {
		file := ts.MkAbs(arg)
		removeAll(file)              // does chmod and then attempts rm
		ts.Check(os.RemoveAll(file)) // report error
	}
}

// symlink creates a symbolic link.
func (ts *Script) CmdSymlink(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? symlink")
	}
	if len(args) != 3 || args[1] != "->" {
		ts.Fatalf("usage: symlink file -> target")
	}
	// Note that the link target args[2] is not interpreted with MkAbs:
	// it will be interpreted relative to the directory file is in.
	ts.Check(os.Symlink(args[2], ts.MkAbs(args[0])))
}
