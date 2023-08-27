package runtime

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mattn/go-zglob"

	"github.com/hofstadter-io/hof/lib/hof"
)

func keepFilter(hn *hof.Node[any], patterns []string) bool {
	// filter by name
	if len(patterns) > 0 {
		for _, d := range patterns {

			// three match variations
			// 1. regexp when /.../
			// 2. glob if any *
			// 3. string prefix
			if strings.HasPrefix(d,"/") && strings.HasSuffix(d,"/") {
				// regexp
				match, err := regexp.MatchString(d, hn.Hof.Metadata.Name)
				if err != nil {
					fmt.Println("error:", err)
					return false
				}
				if match {
					return true
				}
			} else if strings.Contains(d,"*") {
				// glob
				match, err := zglob.Match(d, hn.Hof.Metadata.Name)
				if err != nil {
					fmt.Println("error:", err)
					return false
				}
				if match {
					return true
				}
			} else {
				// prefix
				if strings.HasPrefix(hn.Hof.Metadata.Name, d) {
					return true	
				}
			}
		}
		return false
	}

	// filter by time

	// filter by version?

	// default to true, should include everything when no checks are needed
	return true
}

