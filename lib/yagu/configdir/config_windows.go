package configdir

import "os"

var (
	systemConfig []string
	localConfig  string
	localCache   string
)

func findPaths() {
	systemConfig = []string{os.Getenv("PROGRAMDATA")}
	localConfig = os.Getenv("APPDATA")
	localCache = os.Getenv("LOCALAPPDATA")
}
