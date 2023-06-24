package dagger

import (
	"fmt"

	"dagger.io/dagger"
)


func (R *Runtime) fmtrContainer(name, tag string) (*dagger.Container, error) {
		
	c := R.Client.Container().From(fmt.Sprintf("ghcr.io/hofstadter-io/fmt-%s:%s"))
	c = c.Pipeline(fmt.Sprintf("image/fmt/%s", name))
	c = c.WithExposedPort(3000)

	return c, nil
}

func (R *Runtime) attachFmtr(cntr *dagger.Container, name string, fmtr *dagger.Container) (*dagger.Container, error) {
	cntr = cntr.WithServiceBinding(fmt.Sprintf("hof-fmt-%s", name), fmtr)
	return cntr, nil
}
