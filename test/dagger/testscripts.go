package dagger

import (
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

func (R *Runtime) AddDevTools(c *dagger.Container) (*dagger.Container) {
	c = c.WithExec([]string{"apt-get", "update", "-y"})
	c = c.WithExec([]string{"apt-get", "install", "-y", "tree"})
	return c
}

func (R *Runtime) ShowWorkdir(c *dagger.Container) *dagger.Container {
	c = c.WithExec([]string{"pwd"})
	c = c.WithExec([]string{"tree"})
	return c
}

func (R *Runtime) SetupTestingEnv(c *dagger.Container) (*dagger.Container) {
	c = c.Pipeline("setup/testenv")
	c = c.WithEnvVariable("HOF_TELEMETRY_DISABLED", "1")
	c = c.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	c = c.WithEnvVariable("HOF_FMT_VERSION", os.Getenv("HOF_FMT_VERSION"))
	c = c.WithWorkdir("/test")

	// c = R.AddDevTools(c)
	return c
}

func (R *Runtime) RunTestscriptDir(c *dagger.Container, source *dagger.Directory, name, dir string) error {

	d := source.Directory(dir)
	files, err := d.Entries(R.Ctx)
	if err != nil {
		return err
	}

	p := c.Pipeline(name)

	// we want to run each as a separate fork of the testing container, in this way
	// each test gets a fresh environment and we can collect multiple errors before failing totally
	hadError := false
	for _, f := range files {
		ext := filepath.Ext(f)
		if ext == ".txt" {
			t := p.Pipeline(filepath.Join(dir, f))

			t = t.WithFile(filepath.Join("/test", f), d.File(f))
			t = t.WithExec([]string{"hof", "run", f})

			// now we only sync and check results once
			_, err = t.Sync(R.Ctx)
			if err != nil {
				hadError = true
			}
		}
	}

	if hadError {
		return fmt.Errorf("errors while running %s in %s", name, dir)
	}

	return nil
}

func (R *Runtime) TestCommandFmt(c *dagger.Container, source *dagger.Directory) error {

	t := c.Pipeline("test/fmt")

	ver := "v0.6.8-rc.5"

	daemon, err := R.DockerDaemonContainer()
	if err != nil {
		return err
	}

	t, err = R.AttachDaemonAsService(t, daemon)
	if err != nil {
		return err
	}

	t = t.WithEnvVariable("HOF_FMT_HOST", "http://global-dockerd")
	t = t.WithEnvVariable("HOF_FMT_VERSION", ver)

	t = t.WithExec([]string{"hof", "fmt", "start", "all"})
	t = t.WithExec([]string{"hof", "fmt", "info"})

	err = R.RunTestscriptDir(t, source, "test/fmt", "formatters/test")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestAdhocRender(c *dagger.Container, source *dagger.Directory) error {
	t := c.Pipeline("test/render")

	err := R.RunTestscriptDir(t, source, "test/render", "test/render")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestMod(c *dagger.Container, source *dagger.Directory) error {
	t := c.Pipeline("test/mod")

	t = t.WithEnvVariable("GITLAB_TOKEN", os.Getenv("GITLAB_TOKEN"))
	t = t.WithEnvVariable("BITBUCKET_USERNAME", os.Getenv("BITBUCKET_USERNAME"))
	t = t.WithEnvVariable("BITBUCKET_PASSWORD", os.Getenv("BITBUCKET_PASSWORD"))

	err := R.RunTestscriptDir(t, source, "test/mod", "lib/mod/testdata")
	if err != nil {
		return err
	}

	t = t.Pipeline("test/mod/auth")
	err = R.RunTestscriptDir(t, source, "test/mod/auth", "lib/mod/testdata/authd/apikeys")
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) TestCreate(c *dagger.Container, source *dagger.Directory) error {
	p := c.Pipeline("test/create")

	dirs := []string{
		"test/create/test_01",
		"test/create/test_02",
	}

	for _, dir := range dirs {
		d := source.Directory(dir)

		t := p.Pipeline(dir)
		t = t.WithDirectory("/test", d)
		t = t.WithExec([]string{"make", "test"})
		_, err := t.Sync(R.Ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

