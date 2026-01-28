package runtime

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hofstadter-io/hof/lib/config"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/hofstadter-io/hof/lib/container"
)

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
	name := consts.VEG_REGISTRY_DEFAULT_CONTAINER_NAME
	containers, err := container.GetContainers(name)
	if err != nil {
		return err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == "/"+name || n == name {
				if c.Image != consts.VEG_REGISTRY_DEFAULT_IMAGE {
					fmt.Printf("A new version of %s is available: %s (current: %s)\n", name, consts.VEG_REGISTRY_DEFAULT_IMAGE, c.Image)
				}
				if c.State == "running" {
					return nil
				}
			}
		}
	}

	fmt.Println("Starting veg-registry...")

	params := &container.Params{
		Name:    container.Name(name),
		Replace: true,
		Publish: []string{"5000:5000"},
		Restart: "always",
		Volume:  []string{fmt.Sprintf("%s:/var/lib/registry", config.Veg.RegistryData)},
		Env: []string{
			"OTEL_TRACES_EXPORTER=none",
			// "REGISTRY_HTTP_SECRET=vegreg",
		},
	}

	return container.StartContainer(consts.VEG_REGISTRY_DEFAULT_IMAGE, params)
}

func ensureDagger() error {
	name := consts.VEG_DAGGER_ENGINE_DEFAULT_CONTAINER_NAME
	containers, err := container.GetContainers(name)
	if err != nil {
		return err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == "/"+name || n == name {
				if c.Image != consts.VEG_DAGGER_ENGINE_DEFAULT_IMAGE {
					fmt.Printf("A new version of %s is available: %s (current: %s)\n", name, consts.VEG_DAGGER_ENGINE_DEFAULT_IMAGE, c.Image)
				}
				if c.State == "running" {
					return nil
				}
			}
		}
	}

	fmt.Printf("Starting %s...\n", name)

	// use existing config if it exists
	enginePath := config.Veg.DaggerEngineConfig
	fmt.Println("looking for config:", enginePath)
	if _, err := os.Stat(enginePath); os.IsNotExist(err) {
		fmt.Println("creating first engine config:", enginePath)
		err = os.MkdirAll(filepath.Dir(enginePath), 0755)
		if err != nil {
			return err
		}
		err = os.WriteFile(enginePath, []byte(consts.VEG_DAGGER_ENGINE_CONFIG_DEFAULT), 0644)
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
			fmt.Sprintf("%s:/etc/dagger/engine.json", enginePath),
		},
		AddHost: []string{"host.docker.internal:host-gateway"},
	}

	return container.StartContainer(consts.VEG_DAGGER_ENGINE_DEFAULT_IMAGE, params)
}
