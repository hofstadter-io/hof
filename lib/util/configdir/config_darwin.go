package configdir

import "os"

var (
	systemConfig []string
	localConfig  string
	localCache   string
)

func findPaths() {
	systemConfig = []string{"/Library/Application Support"}
	localConfig = os.Getenv("HOME") + "/Library/Application Support"
	localCache = os.Getenv("HOME") + "/Library/Caches"
}
