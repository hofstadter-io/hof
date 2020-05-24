package configdir_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hofstadter-io/hof/lib/yagu/configdir"
)

// Quick start example for a common use case of this module.
func Example() {
	// A common use case is to get a private config folder for your app to
	// place its settings files into, that are specific to the local user.
	configPath := configdir.LocalConfig("my-app")
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		panic(err)
	}

	// Deal with a JSON configuration file in that folder.
	configFile := filepath.Join(configPath, "settings.json")
	type AppSettings struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var settings AppSettings

	// Does the file not exist?
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		// Create the new config file.
		settings = AppSettings{"MyUser", "MyPassword"}
		fh, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		encoder := json.NewEncoder(fh)
		encoder.Encode(&settings)
	} else {
		// Load the existing file.
		fh, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		decoder := json.NewDecoder(fh)
		decoder.Decode(&settings)
	}
}

// Example for getting a global system configuration path.
func ExampleSystemConfig() {
	// Get all of the global system configuration paths.
	//
	// On Linux or BSD this might be []string{"/etc/xdg"} or the split
	// version of $XDG_CONFIG_DIRS.
	//
	// On macOS or Windows this will likely return a slice with only one entry
	// that points to the global config path; see the README.md for details.
	paths := configdir.SystemConfig()
	fmt.Printf("Global system config paths: %v\n", paths)

	// Or you can get a version of the path suffixed with a vendor folder.
	vendor := configdir.SystemConfig("acme")
	fmt.Printf("Vendor-specific config paths: %v\n", vendor)

	// Or you can use multiple path suffixes to group configs in a
	// `vendor/application` namespace. You can use as many path
	// components as you like.
	app := configdir.SystemConfig("acme", "sprockets")
	fmt.Printf("Vendor/app specific config paths: %v\n", app)
}

// Example for getting a user-specific configuration path.
func ExampleLocalConfig() {
	// Get the default root of the local configuration path.
	//
	// On Linux or BSD this might be "$HOME/.config", or on Windows this might
	// be "C:\\Users\\$USER\\AppData\\Roaming"
	path := configdir.LocalConfig()
	fmt.Printf("Local user config path: %s\n", path)

	// Or you can get a local config path with a vendor suffix, like
	// "$HOME/.config/acme" on Linux.
	vendor := configdir.LocalConfig("acme")
	fmt.Printf("Vendor-specific local config path: %s\n", vendor)

	// Or you can use multiple path suffixes to group configs in a
	// `vendor/application` namespace. You can use as many path
	// components as you like.
	app := configdir.LocalConfig("acme", "sprockets")
	fmt.Printf("Vendor/app specific local config path: %s\n", app)
}

// Example for getting a user-specific cache folder.
func ExampleLocalCache() {
	// Get the default root of the local cache folder.
	//
	// On Linux or BSD this might be "$HOME/.cache", or on Windows this might
	// be "C:\\Users\\$USER\\AppData\\Local"
	path := configdir.LocalCache()
	fmt.Printf("Local user cache path: %s\n", path)

	// Or you can get a local cache path with a vendor suffix, like
	// "$HOME/.cache/acme" on Linux.
	vendor := configdir.LocalCache("acme")
	fmt.Printf("Vendor-specific local cache path: %s\n", vendor)

	// Or you can use multiple path suffixes to group caches in a
	// `vendor/application` namespace. You can use as many path
	// components as you like.
	app := configdir.LocalCache("acme", "sprockets")
	fmt.Printf("Vendor/app specific local cache path: %s\n", app)
}

// Example for automatically creating config directories.
func ExampleMakePath() {
	// The MakePath() function can accept the output of any of the folder
	// getting functions and ensure that their path exists.

	// Create a local user configuration folder under an app prefix.
	// On Linux this may result in `$HOME/.config/my-cool-app` existing as
	// a directory, depending on the value of `$XDG_CONFIG_HOME`.
	err := configdir.MakePath(configdir.LocalConfig("my-cool-app"))
	if err != nil {
		panic(err)
	}

	// Create a cache folder under a namespace.
	err = configdir.MakePath(configdir.LocalCache("acme", "sprockets", "client"))
	if err != nil {
		panic(err)
	}

	// In the case of global system configuration, which may return more than
	// one path (especially on Linux/BSD that uses the XDG Base Directory Spec),
	// it will attempt to create the directories only under the *first* path.
	//
	// For example, if $XDG_CONFIG_DIRS="/etc/xdg:/opt/config" this will try
	// to create the config dir only in "/etc/xdg/acme/sprockets" and not try
	// to create any folders under "/opt/config".
	err = configdir.MakePath(configdir.SystemConfig("acme", "sprockets")...)
	if err != nil {
		panic(err)
	}
}

// Example for recalculating what the directories should be.
func ExampleRefresh() {
	// On your program's initialization, this module decides which paths to
	// use for global, local and cache folders, based on environment variables
	// and falling back on defaults.
	//
	// In case the environment variables change throughout the life of your
	// program, for example if you re-assigned $XDG_CONFIG_HOME, you can call
	// the Refresh() function to re-calculate the paths to reflect the new
	// environment.
	configdir.Refresh()
}
