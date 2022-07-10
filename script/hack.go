package script

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/script/_ast"
)

func Hack(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("please supply a single filepath")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fs := osfs.New(cwd)

	llvls := []string{"error", "warn", "info", "debug"}
	llvl := llvls[flags.RootPflags.Verbosity]

	config := &ast.Config{
		LogLevel: llvl,
		FS:       fs,
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
