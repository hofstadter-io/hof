package dagger

import (
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)


// minimal implementation of test script in dagger
func (R *Runtime) RenderTests(c *dagger.Container) error {
	name := "test/render"
	return R.runTestscriptDir(c, name)
}

func (R *Runtime) runTestscriptDir(c *dagger.Container, dir string) error {

	d := c.Directory(dir)
	files, err := d.Entries(R.Ctx)
	if err != nil {
		return err
	}

	p := c.Pipeline(dir)
	p = p.WithEnvVariable("HOF_TELEMETRY_DISABLED", "1")
	p = p.WithEnvVariable("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))
	p = p.WithEnvVariable("HOF_FMT_VERSION", os.Getenv("HOF_FMT_VERSION"))
	p = p.WithWorkdir("/test")

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

func (R *Runtime) runTestscriptFile(c *dagger.Container, filepath string) error {

	return nil
}
