package dotpath

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func get_from_slice_by_path(IDX int, paths []string, data []interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("input data is nil")
	}

	header := fmt.Sprintf("get_from_slice_by_path:  %d  %v  in:\n%+v\n\n", IDX, paths, data)
	// fmt.Println(header)
	logger.Debug(header)

	if IDX >= len(paths) || len(paths) == 0 {
		return nil, nil
	}

	subs := []interface{}{}

	P := paths[IDX]

	lpos_index := strings.Index(P, "[")
	rpos_index := strings.LastIndex(P, "]")
	pos_colon := strings.Index(P, ":")
	has_listing := strings.Contains(P, ",")
	// pos_regex := strings.Index(P, "regex")

	has_eq := strings.Contains(P, "==")
	// has_ne := strings.Contains(P, "!=")

	inner := ""
	if lpos_index > -1 {
		inner = P[lpos_index+1 : rpos_index]
	}
	// fmt.Printf("  slice inner: %d %q %q\n", IDX, inner, P)

	// handle indexing here
	if inner != "" {
		inner := P[lpos_index+1 : rpos_index]
		// fmt.Printf("index: %q  [%d:%d]\n", inner, lpos_index+1, rpos_index)

		// handle slicing
		if pos_colon > -1 {
			elems, err := extract_from_slice_with_splice(inner, data)
			if err != nil {
				return nil, errors.Wrap(err, "while extracting splice in dotpath case []interface{}")
			}
			return elems, nil
		}

		// handle equality

		if has_eq {
			fields := strings.Split(inner, "==")
			if len(fields) != 2 {
				return nil, errors.New("Found not 2 fields in equality in: " + P)
			}
			elems, err := extract_from_slice_with_field(fields[0], fields[1], data)
			if err != nil {
				return nil, errors.Wrap(err, "while extracting has_eq in dotpath case []interface{}")
			}
			return elems, nil
		}

		// handle listing
		if has_listing {
			elems, err := extract_from_slice_with_name(inner, data)
			if err != nil {
				return nil, errors.Wrap(err, "while extracting listing in dotpath case []interface{}")
			}
			return elems, nil
		}

		// handle int
		i_val, ierr := strconv.Atoi(inner)
		if ierr == nil {
			if i_val < 0 || i_val >= len(data) {
				str := "index out of range: " + inner
				return str, errors.New(str)
			}
			return data[i_val], nil
		}

		// default is single string for name, then field
		elems, err := extract_from_slice_with_name(inner, data)
		if err != nil {
			return nil, errors.Wrap(err, "while extracting name/field in dotpath case []interface{}")
		}

		return elems, nil

	} else {
		// No inner indexing
		for _, elem := range data {
			logger.Debug("    - elem", "elem", elem, "paths", paths, "P", P)
			switch V := elem.(type) {

			case map[string]interface{}:
				logger.Debug("        map[string]")
				val, err := get_from_smap_by_path(IDX, paths, V)
				if err != nil {
					logger.Debug("could not find '" + P + "' in object")
					continue
				}
				logger.Debug("Adding val", "val", val)
				subs = append(subs, val)

			case map[interface{}]interface{}:
				logger.Debug("        map[iface]")
				val, err := get_from_imap_by_path(IDX, paths, V)
				if err != nil {
					logger.Debug("could not find '" + P + "' in object")
					continue
				}
				logger.Debug("Adding val", "val", val)
				subs = append(subs, val)

			default:
				str := fmt.Sprintf("%+v", reflect.TypeOf(V))
				return nil, errors.New("element not an object type: " + str)

			}
		}
	}

	if len(subs) == 1 {
		return subs[0], nil
	}
	return subs, nil

}
