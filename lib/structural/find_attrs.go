package structural

import (
	"cuelang.org/go/cue"
)

var xx struct{}

// FindByAttrs finds values in s if
// 1) attrs : it has @attr(...)
// 2) attrsK : it has @attr(key,...) [it must have all keys in the passed list, not just one]
// 3) attrsKV : it has @attr(key=value,...) [it must have all key=value in the passed map, not just one of them]
func FindByAttrs(s cue.Value, attrs []string, attrsK map[string][]string, attrsKV map[string]map[string]string) ([]cue.Value, error) {
	out := make([]cue.Value, 0)

	siter, err := s.Fields()
	if err != nil {
		return nil, err
	}

	attrsSet := make(map[string]struct{})
	for _, a := range attrs {
		attrsSet[a] = xx
	}

	for siter.Next() {
		// label := siter.Label()
		value := siter.Value()
		attrs := value.Attributes()
		for _, attr := range attrs {
			k := attr.Name()
			if _, ok := attrsSet[k]; ok {
				out = append(out, value)
				break
			}
			if ks, ok := attrsK[k]; ok {
				include := true
				for _, checkKey := range ks {
					// TODO API is lacking here, assume we have less than 20 attribute val positions...
					found := false
					for i := 0; i < 20; i++ {
						key, err := attr.String(i)
						if err == nil && key == checkKey {
							found = true
							break
						}
					}
					if !found {
						include = false
						break
					}
				}
				if include {
					out = append(out, value)
					break
				}
			}
			if kvs, ok := attrsKV[k]; ok {
				include := true
				for checkKey, checkVal := range kvs {
					val, found, err := attr.Lookup(0, checkKey)
					if !found || err != nil || val != checkVal {
						include = false
						break
					}
				}
				if include {
					out = append(out, value)
					break
				}
			}
		}
	}

	return out, nil
}
