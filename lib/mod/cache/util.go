package cache

import "strings"

func parseModURL(mod string) (remote, owner, repo string) {
	var flds []string
	if i := strings.Index(mod, ".git"); i > -1 {
		flds = strings.SplitN(mod[:i], "/", 3)
	} else {
		flds = strings.Split(mod, "/")
	}

	return flds[0], flds[1], flds[2]
}
