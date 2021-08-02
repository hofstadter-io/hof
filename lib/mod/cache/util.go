package cache

import (
	"path"
	"strings"
)

func parseModURL(mod string) (remote, owner, repo string) {
	var flds []string
	if i := strings.Index(mod, ".git"); i > -1 {
		flds = strings.SplitN(mod[:i], "/", 3)
	} else {
		flds = strings.Split(mod, "/")
	}

	if len(flds) < 3 {
		return flds[0], "", flds[1]
	}

	return flds[0], flds[1], flds[2]
}

// MatchPrefixPatterns reports whether any path prefix of target matches one of
// the glob patterns (as defined by path.Match) in the comma-separated globs
// list. This implements the algorithm used when matching a module path to the
// GOPRIVATE environment variable, as described by 'go help module-private'.
//
// It ignores any empty or malformed patterns in the list.
func MatchPrefixPatterns(globs, target string) bool {
	for globs != "" {
		// Extract next non-empty glob in comma-separated list.
		var glob string
		if i := strings.Index(globs, ","); i >= 0 {
			glob, globs = globs[:i], globs[i+1:]
		} else {
			glob, globs = globs, ""
		}
		if glob == "" {
			continue
		}

		// A glob with N+1 path elements (N slashes) needs to be matched
		// against the first N+1 path elements of target,
		// which end just before the N+1'th slash.
		n := strings.Count(glob, "/")
		prefix := target
		// Walk target, counting slashes, truncating at the N+1'th slash.
		for i := 0; i < len(target); i++ {
			if target[i] == '/' {
				if n == 0 {
					prefix = target[:i]
					break
				}
				n--
			}
		}
		if n > 0 {
			// Not enough prefix elements.
			continue
		}
		matched, _ := path.Match(glob, prefix)
		if matched {
			return true
		}
	}
	return false
}
