package envs

import (
	"dagger.io/dagger"
)

func (R *Runtime) HarmonizeCue(c *dagger.Container, source *dagger.Directory) (*dagger.Directory) {
	// mount source
	c = c.WithDirectory("/repo", source)	
	c = c.WithWorkdir("/repo")

	// replace CUE version
	c = c.WithExec([]string{ "go", "mod", "edit", "-replace", "cuelang.org/go=cuelang.org/go@" + R.CueVer })

	return c.Directory("/repo")
}

func (R *Runtime) HarmonizeHof(c *dagger.Container, source *dagger.Directory) (*dagger.Directory) {
	c = c.WithDirectory("/repo", source)	
	c = c.WithWorkdir("/repo")

	// replace Hof version
	c = c.WithExec([]string{ "go", "mod", "edit", "-replace", "github.com/hofstadter-io/hof=/work" })

	return c.Directory("/repo")
}
