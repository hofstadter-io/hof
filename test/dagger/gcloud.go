package dagger

import (
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

func (R *Runtime) GcloudImage() (*dagger.Container) {

	c := R.Client.Container().From("google/cloud-sdk")
	c = c.Pipeline("gcloud-sdk")

	return c
}

func (R* Runtime) WithLocalGcloudConfig(c *dagger.Container) (*dagger.Container) {
	cfg, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	d := R.Client.Host().Directory(filepath.Join(cfg, "gcloud"))
	c = c.WithEnvVariable("CLOUDSDK_CONFIG", "/gcloud/config")
	c = c.WithMountedDirectory("/gcloud/config", d)

	return c
}

func (R* Runtime) WithLocalSSHDir(c *dagger.Container) (*dagger.Container) {
	cfg, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	d := R.Client.Host().Directory(filepath.Join(cfg, ".ssh"))
	c = c.WithMountedDirectory("/root/.ssh", d)

	return c
}
