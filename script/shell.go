package script

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/script/ast"
)

func Shell(args []string) error {
	if len(args) > 0 {
		if len(args) != 1 {
			return fmt.Errorf("please supply a single filepath to preload with")
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fs := osfs.New(cwd)

	llvl := "warn"
	if flags.RootPflags.Verbose != "" {
		llvl = flags.RootPflags.Verbose
	}

	config := &ast.Config{
		LogLevel: llvl,
		FS: fs,
	}
	parser := ast.NewParser(config)

	S, err := parser.ParseScript(args[0])
	if err != nil {
		fmt.Println("ERROR:", err)
		return err
	}

	fmt.Println("done hacking ", S.Path)

	return nil
}
