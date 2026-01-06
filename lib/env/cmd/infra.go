package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hofstadter-io/hof/lib/container"
)

const (
	RegistryImage = "registry:3"
	DaggerImage   = "registry.dagger.io/engine:v0.19.8"
)

const DaggerEngineConfig = `
{
  "registries": {
    "host.docker.internal:5000": {
      "http": true
    }
  },
  "gc": {
    "enabled": true,
    "reservedSpace": "80GB",
    "maxUsedSpace": "100GB",
    "minFreeSpace": "10GB"
  }
}
`

func EnsureInfra() error {
	err := container.InitClient()
	if err != nil {
		return err
	}

	err = ensureRegistry()
	if err != nil {
		return err
	}

	err = ensureDagger()
	if err != nil {
		return err
	}

	return nil
}

func ensureRegistry() error {
	name := "veg-registry"
	containers, err := container.GetContainers(name)
	if err != nil {
		return err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == "/"+name || n == name {
				if c.Image != RegistryImage {
					fmt.Printf("A new version of %s is available: %s (current: %s)\n", name, RegistryImage, c.Image)
				}
				if c.State == "running" {
					return nil
				}
			}
		}
	}

	fmt.Println("Starting veg-registry...")

	registryData := os.Getenv("VEG_REGISTRY_DATA")
	if registryData == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return err
		}
		registryData = filepath.Join(cacheDir, "veg", "registry")
	}

	params := &container.Params{
		Name:    container.Name(name),
		Replace: true,
		Publish: []string{"5000:5000"},
		Restart: "always",
		Volume:  []string{fmt.Sprintf("%s:/var/lib/registry", registryData)},
	}

	return container.StartContainer(RegistryImage, params)
}

func ensureDagger() error {
	name := "veg-dagger-engine"
	containers, err := container.GetContainers(name)
	if err != nil {
		return err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == "/"+name || n == name {
				if c.Image != DaggerImage {
					fmt.Printf("A new version of %s is available: %s (current: %s)\n", name, DaggerImage, c.Image)
				}
				if c.State == "running" {
					return nil
				}
			}
		}
	}

	fmt.Println("Starting veg-dagger-engine...")

	// use existing config if it exists
	enginePath := os.Getenv("VEG_DAGGER_ENGINE_JSON")
	if enginePath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		enginePath = filepath.Join(configDir, "veg", "dagger-engine.json")
	}

	fmt.Println("looking for config:", enginePath)
	if _, err := os.Stat(enginePath); os.IsNotExist(err) {
		fmt.Println("creating first engine config:", enginePath)
		err = os.MkdirAll(filepath.Dir(enginePath), 0755)
		if err != nil {
			return err
		}
		err = os.WriteFile(enginePath, []byte(DaggerEngineConfig), 0644)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("found config at:", enginePath)
	}

	params := &container.Params{
		Name:       container.Name(name),
		Replace:    true,
		Restart:    "always",
		Privileged: true,
		Volume: []string{
			"/var/lib/dagger",
			// fmt.Sprintf("%s:/etc/dagger/engine.json", enginePath),
		},
		AddHost: []string{"host.docker.internal:host-gateway"},
	}

	return container.StartContainer(DaggerImage, params)
}
