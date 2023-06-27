package dagger

import (
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

func (R *Runtime) GcloudImage() (*dagger.Container) {

	cfg, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	d := R.Client.Host().Directory(filepath.Join(cfg, "gcloud"))

	c := R.Client.Container().From("google/cloud-sdk")
	c = c.Pipeline("gcloud-sdk")
	c = c.WithEnvVariable("CLOUDSDK_CONFIG", "/gcloud/config")
	c = c.WithDirectory("/gcloud/config", d)

	return c
}
