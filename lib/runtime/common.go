package runtime

import (
	"fmt"
	"regexp"

	"github.com/hofstadter-io/hof/lib/hof"
)

func keepFilter(hn *hof.Node[any], patterns []string) bool {
	// filter by name
	if len(patterns) > 0 {
		for _, d := range patterns {
			match, err := regexp.MatchString(d, hn.Hof.Metadata.Name)
			if err != nil {
				fmt.Println("error:", err)
				return false
			}
			if match {
				return true
			}
		}
		return false
	}

	// filter by time

	// filter by version?

	// default to true, should include everything when no checks are needed
	return true
}

