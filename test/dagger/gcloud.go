package dagger

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dagger.io/dagger"
)

func (R *Runtime) GcloudImage() (*dagger.Container) {

	c := R.Client.Container().From("google/cloud-sdk")
	c = c.Pipeline("gcloud-sdk")

	return c
}

func (R* Runtime) WithLocalGcloudConfig(c *dagger.Container) (*dagger.Container) {
	out, err := exec.Command("bash", "-c", "gcloud info --format='value(config. paths. global_config_dir)'").CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	cfg := strings.TrimSpace(string(out))
	d := R.Client.Host().Directory(cfg)
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
