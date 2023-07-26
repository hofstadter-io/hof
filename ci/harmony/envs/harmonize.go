package envs

import (
	"fmt"

	"dagger.io/dagger"
)

func (R *Runtime) HarmonizeCue(c *dagger.Container, source *dagger.Directory) (*dagger.Directory) {
	// mount source
	c = c.WithDirectory("/work", source)	

	// replace CUE version
	c = c.WithExec([]string{
		"bash",
		"-c",
		"go mod edit -replace cuelang.org/go=cuelang.org/go@" + R.CueVer,
	})

	return c.Directory("/work")
}

func (R *Runtime) HarmonizeHof(c *dagger.Container, source *dagger.Directory) (*dagger.Directory) {
	c = c.WithDirectory("/repo", source)	

	repo := ""

	// replace CUE version
	c = c.WithExec([]string{
		"bash",
		"-c",
		fmt.Sprintf("go mod edit -replace github.com/hofstadter-io/hof=/work", repo, repo, R.HofVer),
	})

	return c.Directory("/work")
}
