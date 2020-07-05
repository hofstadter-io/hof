package runtime

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hofstadter-io/hof/script/ast"
)

func (RT *Runtime) Cmd_pwd(cmd *ast.Cmd, r *ast.Result) (err error) {

	if cmd.Exp != ast.Pass {
		return fmt.Errorf("unsupported: !? pwd")
	}

	cwd := RT.currdir

	fmt.Fprintln(cmd.Result().Stdout, cwd)

	return nil
}

func (RT *Runtime) Cmd_mkdir(cmd *ast.Cmd, r *ast.Result) (err error) {

	cwd := RT.currdir

	fmt.Fprintln(cmd.Result().Stdout, cwd)

	return nil
}

func (RT *Runtime) Cmd_cd(cmd *ast.Cmd, r *ast.Result) (err error) {
	RT.nextCmd(cmd)

	cwd := RT.currdir

	if len(cmd.Args) == 0 {
		// go home
		cwd, err = os.UserHomeDir()
		if err != nil {
			r.AddError(err)
			return err
		}
	} else if len(cmd.Args) == 1 {
		cwd = filepath.Join(cwd, cmd.Args[0])
		cwd, err = filepath.Rel("/", cwd)
		if err != nil {
			return err
		}
		cwd = filepath.Join("/", cwd)
	} else {
		return fmt.Errorf("Too many args to cd")
	}

	RT.currdir = cwd

	return nil
}

func (RT *Runtime) Cmd_ls(cmd *ast.Cmd, r *ast.Result) (err error) {

	args := cmd.Args
	if len(args) == 0 {
		args = []string{"."}
	}

	for _, dir := range args {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			// ok here, as we are building up errors on ls args
			r.AddError(err)
			r.Status = 1
			RT.status = r.Status
			continue
		}

		for _, file := range files {
			RT.printFileInfo(cmd.Result().Stdout, file)
		}
	}

	if r.Status < 0 {
		r.Status = 0
		RT.status = r.Status
	}

	return  nil
}

// MkAbs interprets file relative to the test script's current directory
// and returns the corresponding absolute path.
func (RT *Runtime) MkAbs(file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	return filepath.Join(RT.currdir, file)
}


func (RT *Runtime) printFileInfo(w io.Writer, file os.FileInfo) {
		fmt.Fprintf(w, "%s\n", file.Name())
}


// ReadFile returns the contents of the file with the
// given name, intepreted relative to the test script's
// current directory. It interprets "stdout" and "stderr" to
// mean the standard output or standard error from
// the most recent exec or wait command respectively.
//
// If the file cannot be read, the script fails.
func (RT *Runtime) ReadFile(file string) string {
	switch file {
	case "stdout":
		return RT.stdoutStr
	case "stderr":
		return RT.stderrStr
	default:
		file = RT.MkAbs(file)
		data, err := ioutil.ReadFile(file)
		RT.Check(err)
		return string(data)
	}
}

func removeAll(dir string) error {
	// module cache has 0444 directories;
	// make them writable in order to remove content.
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // ignore errors walking in file system
		}
		if info.IsDir() {
			os.Chmod(path, 0777)
		}
		return nil
	})
	return os.RemoveAll(dir)
}
