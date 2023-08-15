package verinfo

import (
	"runtime"
	"runtime/debug"
)

var (
	Version = "Local"
	Commit  = "Dirty"

	BuildDate = "Unknown"
	GoVersion = "Unknown"
	BuildOS   = "Unknown"
	BuildArch = "Unknown"

	// todo, look this up from deps
	CueVersion = "0.6.0"

	// this is a version we can fetch with hof mod
	// the value gets injected into templates in various places
	// the default here is set to something useful for dev
	// the release version is the same as the cli running it
	HofVersion = "v0.6.8"
)


func init() {
	info, _ := debug.ReadBuildInfo()
	GoVersion = info.GoVersion

	if Version == "Local" {
		BuildOS = runtime.GOOS
		BuildArch = runtime.GOARCH

		dirty := false
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				Commit = s.Value
			}
			if s.Key == "vcs.time" {
				BuildDate = s.Value
			}
			if s.Key == "vcs.modified" {
				if s.Value == "true" {
					dirty = true
				}
			}
		}
		if dirty {
			Commit += "+dirty"
		}
	}

	// released binary override
	if Version != "Local" {
		HofVersion = Version
	}
}
