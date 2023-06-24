package dagger

import (
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

func (R *Runtime) SetupTestingEnv(c *dagger.Container) (*dagger.Container) {
	c = c.WithEnvVariable("HOF_TELEMETRY_DISABLED", "1")
	c = c.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	c = c.WithEnvVariable("HOF_FMT_VERSION", os.Getenv("HOF_FMT_VERSION"))
	c = c.WithWorkdir("/test")
	return c
}

func (R *Runtime) RunTestscriptDir(c *dagger.Container, source *dagger.Directory, name, dir string) error {

	d := source.Directory(dir)
	files, err := d.Entries(R.Ctx)
	if err != nil {
		return err
	}

	p := c.Pipeline(name)

	errs := []error{}

	for _, f := range files {
		ext := filepath.Ext(f)
		if ext == ".txt" {

			t := p.Pipeline(filepath.Join(dir, f))
			t = t.WithFile(filepath.Join("/test", f), d.File(f))

			t = t.WithExec([]string{"hof", "run", f})
			code, exiterr := t.ExitCode(R.Ctx)
			stdout, outerr := t.Stdout(R.Ctx)
			stderr, errerr := t.Stderr(R.Ctx)
			if exiterr != nil {
				errs = append(errs, fmt.Errorf("while getting exit code: %s\n%s\n%s", stdout, stderr, exiterr))
				continue
			}
			if outerr != nil {
				errs = append(errs, fmt.Errorf("while getting stdout: %s\n%s\n%s", stdout, stderr, outerr))
				continue
			}
			if errerr != nil {
				errs = append(errs, fmt.Errorf("while getting stderr: %s\n%s\n%s", stdout, stderr, errerr))
				continue
			}
			if code != 0 {
				errs = append(errs, fmt.Errorf("%s\n%s\nexited with error %v", stdout, stderr, code))
				continue
			}
		}
	}

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println("Error:", e)
		}
		return fmt.Errorf("%d errors while running %s\n", len(errs), dir)
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
