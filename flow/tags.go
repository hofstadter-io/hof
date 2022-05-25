package flow

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/structural"
)

func getTagsAndSecrets(val cue.Value) (tags []cue.Value, secrets []cue.Value, errs []error) {

	// fuction used during tree walk to collect values with tags
	collector := func(v cue.Value) bool {
		attrs := v.Attributes(cue.ValueAttr)

		var err error
		for _, attr := range attrs {
			if attr.Name() == "tag" {
				if attr.NumArgs() == 0 {
					err = fmt.Errorf("@tag() has no inner args at %s", v.Path())
					errs = append(errs, err)
				}
				tags = append(tags, v)
			}
			if attr.Name() == "secret" {
				secrets = append(secrets, v)
			}
		}

		return true
	}

	structural.Walk(val, collector, nil, walkOptions...)

	return tags, secrets, errs
}

func injectTags(val cue.Value, tags []string) (cue.Value, error) {
	tagMap := make(map[string]string)
	for _, t := range tags {
		fs := strings.SplitN(t, "=", 2)
		if len(fs) != 2 {
			return val, fmt.Errorf("tags must have form key=value, got %q", t)
		}
		tagMap[fs[0]] = fs[1]
	}

	tagPaths := make(map[string]cue.Path)
	errs := []error{}
	collector := func(v cue.Value) bool {
		attrs := v.Attributes(cue.ValueAttr)

		var err error
		for _, attr := range attrs {
			if attr.Name() == "tag" {
				if attr.NumArgs() == 0 {
					err = fmt.Errorf("@tag() has no inner args at %s", v.Path())
					errs = append(errs, err)
					return false
				}
				// TODO, better options &| UX here
				arg, _ := attr.String(0)
				_, ok := tagMap[arg]
				if ok {
					tagPaths[arg] = v.Path()
				}

				return false
			}
		}

		return true
	}

	structural.Walk(val, collector, nil, walkOptions...)

	for arg, path := range tagPaths {
		val = val.FillPath(path, tagMap[arg])
	}

	return val, nil
}
