package util

import (
	"github.com/bmatcuk/doublestar"
)

func CheckShouldInclude(filename string, includes, excludes []string) (bool, error) {
	var err error
	include := false

	if len(includes) > 0 {
		for _, pattern := range includes {
			include, err = doublestar.PathMatch(pattern, filename)
			if err != nil {
				return false, err
			}
			if include {
				break
			}
		}
	} else {
		include = true
	}

	exclude := false
	if len(excludes) > 0 {
		for _, pattern := range excludes {
			exclude, err = doublestar.PathMatch(pattern, filename)
			if err != nil {
				return false, err
			}
			if exclude {
				break
			}
		}
	}

	return include && !exclude, nil
}
