package verinfo

import (
	"runtime"
	"time"
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

	if BuildDate == "Unknown" {
		BuildDate = time.Now().String()
		GoVersion = "run 'go version', you should have been the one who built this"
		BuildOS = runtime.GOOS
		BuildArch = runtime.GOARCH
	}
}
