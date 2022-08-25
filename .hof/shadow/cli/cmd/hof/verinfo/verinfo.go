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
}
