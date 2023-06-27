package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"dagger.io/dagger"
	hdagger "github.com/hofstadter-io/hof/test/dagger"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	R := &hdagger.Runtime{
		Ctx:    ctx,
		Client: client,
	}

	// docker := R.DockerImage()
	gcloud := R.GcloudImage()

	t := gcloud
	t = t.WithEnvVariable("CACHEBUST", time.Now().String())
	t = t.Pipeline("gcloud/list")
	t = t.WithExec([]string{"gcloud", "compute", "instances", "list", "--format=json"})

	out, err := t.Stdout(ctx)
	checkErr(err)

	vals := make([]map[string]any, 0)
	err = json.Unmarshal([]byte(out), &vals)
	checkErr(err)

	for _, val := range vals {
		name := val["name"].(string)
		zone := val["zone"].(string)

		// skip k8s vms
		if strings.Contains(name, "gke") {
			continue
		}

		d := t.Pipeline("gcloud/describe/" + name)
		d = d.WithExec([]string{"gcloud", "compute", "instances", "describe", name, "--format=json", "--zone", zone})
		d.Sync(ctx)
	}


	final := t
	final.Sync(ctx)
	// final.Stdout(ctx)
}
