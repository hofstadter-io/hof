package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hofstadter-io/hof/lib/consts"
)

func init() {
	Veg = new(VegConfig)

	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("WARN: unable to find user config dir:", err)
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("WARN: unable to find user cache dir:", err)
	}

	// Directory Settings
	// we don't actually use this yet, so let's not every set it either
	// Veg.SysDataDir = os.Getenv(consts.VEG_SYS_DATA_DIR_VAR)
	// if Veg.SysDataDir == "" {
	// 	Veg.SysDataDir = consts.VEG_SYS_DATA_DIR_VAR
	// }
	Veg.UserDataDir = os.Getenv(consts.VEG_USER_DATA_DIR_VAR)
	if Veg.UserDataDir == "" {
		Veg.UserDataDir = filepath.Join(configDir, "veg", "data")
	}

	// Database Settings
	Veg.DatabaseType = os.Getenv(consts.VEG_DATABASE_TYPE_VAR)
	if Veg.DatabaseType == "" {
		Veg.DatabaseType = consts.VEG_DATABASE_TYPE_DEFAULT
	}
	Veg.DatabaseConn = os.Getenv(consts.VEG_DATABASE_CONN_VAR)
	if Veg.DatabaseConn == "" {
		Veg.DatabaseConn = filepath.Join(Veg.UserDataDir, consts.VEG_DATABASE_CONN_DEFAULT)
	}

	// Registry Settings
	Veg.RegistryData = os.Getenv(consts.VEG_REGISTRY_DATA_VAR)
	if Veg.RegistryData == "" {
		Veg.RegistryData = filepath.Join(cacheDir, "veg", "data", "registry")
	}

	// Dagger Settings
	Veg.DaggerEngineConfig = os.Getenv(consts.VEG_DAGGER_ENGINE_CONFIG_VAR)
	if Veg.DaggerEngineConfig == "" {
		// todo, look for multiple formats
		Veg.DaggerEngineConfig = filepath.Join(configDir, "veg", "dagger-engine.json")
	}

	os.MkdirAll(filepath.Join(configDir, "veg"), 0755)
	os.MkdirAll(filepath.Join(cacheDir, "veg"), 0755)
	os.MkdirAll(Veg.UserDataDir, 0755)
	os.MkdirAll(Veg.RegistryData, 0755)

}
