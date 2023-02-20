package main

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/docker"
)

func main() {
	fmt.Println("testing docker client/server compat")


	err := docker.InitDockerClient()
	if err != nil {
		fmt.Println("Error Initializing Client:", err)
		return
	}

	client, server, err := docker.GetVersion()
	if err != nil {
		fmt.Println("Error getting versions:", err)
		return
	}

	fmt.Println("client: ", client)
	fmt.Println("server:", server.Version, server.APIVersion, server.MinAPIVersion)

	img := "hofstadter/fmt-black:v0.6.8-beta.1"

	err = docker.PullImage(img)
	if err != nil {
		fmt.Println("Error pulling image:", err)
		return
	}

	images, err := docker.GetImages(img)
	if err != nil {
		fmt.Println("Error detailing image:", err)
		return
	}

	for _, image := range images {
		fmt.Println(image.ID, image.RepoTags)
	}

}
