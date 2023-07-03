package main

import (
	"fmt"
	"os"

	"github.com/hofstadter-io/hof/lib/container"
)

func main() {
	fmt.Println("testing docker client/server compat")

	err := container.InitClient()
	if err != nil {
		fmt.Println("Error Initializing Client:", err)
		os.Exit(1)
	}

	ver, err := container.GetVersion()
	if err != nil {
		fmt.Println("Error getting versions:", err)
		os.Exit(1)
	}

	fmt.Println("client: ", ver.Client.Version)
	fmt.Println("server:", ver.Server.Version, ver.Server.APIVersion, ver.Server.MinAPIVersion)

	img := "ghcr.io/hofstadter-io/fmt-black:v0.6.8-beta.11"

	err = container.PullImage(img)
	if err != nil {
		fmt.Println("Error pulling image:", err)
		os.Exit(1)
	}

	images, err := container.GetImages(img)
	if err != nil {
		fmt.Println("Error detailing image:", err)
		os.Exit(1)
	}

	for _, image := range images {
		fmt.Println(image.ID, image.RepoTags)
	}
}
