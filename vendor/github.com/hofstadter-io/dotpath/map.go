package dotpath

import (
	"strings"

	"github.com/pkg/errors"
)

func get_from_smap_by_path(IDX int, paths []string, data map[string]interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("input data is nil")
	}

	P := paths[IDX]
	path_str := strings.Join(paths[:IDX+1], ".")

	lpos_index := strings.Index(P, "[")
	rpos_index := strings.LastIndex(P, "]")
	has_listing := strings.Contains(P, ",")
	// pos_regex := strings.Index(P, "regex")

	has_eq := strings.Contains(P, "==")

	inner := ""
	if lpos_index > -1 {
		inner = P[lpos_index+1 : rpos_index]
	}
	// fmt.Printf("  smap inner: %d %q %q\n", IDX, inner, P)

	// handle map indexing
	if inner != "" {
		// handle equality
		if has_eq {
			fields := strings.Split(inner, "==")
			if len(fields) != 2 {
				return nil, errors.New("Found not 2 fields in equality in: " + P)
			}
			elems, err := extract_from_map_by_value(fields[0], fields[1], data)
			if err != nil {
				return nil, errors.Wrap(err, "while extracting has_eq in dotpath case []interface{}")
			}
			return elems, nil
		}

		// handle listing
		if has_listing {
			elems, err := extract_from_map_by_field(inner, data)
			if err != nil {
				return nil, errors.Wrap(err, "while extracting listing in dotpath case []interface{}")
			}
			return elems, nil
		}

		// handle single string field/name
		val, ok := data[inner]
		if !ok {
			// try to look up by name
			name_value, ok := data["name"]
			if ok && name_value == inner {
				return data, nil
			}
			return nil, errors.New("could not find '" + inner + "' in object")
		}
		add_parent_and_path(val, data, path_str)
		if len(paths) == IDX+1 {
			return val, nil
		}
		ret, err := get_by_path(IDX+1, paths, val)
		if err != nil {
			return nil, errors.Wrapf(err, "from object "+inner)
		}
		return ret, nil
	}

	// handle default of field/name
	val, ok := data[P]
	if !ok {
		// try to look up by name
		name_value, ok := data["name"]
		if ok && name_value == P {
			return data, nil
		}
		return nil, errors.New("could not find '" + P + "' in object")
	}
	add_parent_and_path(val, data, path_str)
	if len(paths) == IDX+1 {
		return val, nil
	}
	ret, err := get_by_path(IDX+1, paths, val)
	if err != nil {
		return nil, errors.Wrapf(err, "from object "+P)
	}
	return ret, nil
}

func get_from_imap_by_path(IDX int, paths []string, data map[interface{}]interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("input data is nil")
	}


	P := paths[IDX]
	path_str := strings.Join(paths[:IDX+1], ".")

	lpos_index := strings.Index(P, "[")
	rpos_index := strings.LastIndex(P, "]")
	has_listing := strings.Contains(P, ",")
	// pos_regex := strings.Index(P, "regex")

	has_eq := strings.Contains(P, "==")
	// has_ne := strings.Contains(P, "!=")
	// has_ge := strings.Contains(P, ">=")
	// has_gt := strings.Contains(P, ">")
	// has_le := strings.Contains(P, "<=")
	// has_lt := strings.Contains(P, "<")

	inner := ""
	if lpos_index > -1 {
		inner = P[lpos_index+1 : rpos_index]
	}
	// fmt.Printf("  imap inner: %d %q %q\n", IDX, inner, P)

	// handle map indexing
	if inner != "" {
		// handle equality
		if has_eq {
			fields := strings.Split(inner, "==")
			if len(fields) != 2 {
				return nil, errors.New("Found not 2 fields in equality in: " + P)
			}
			elems, err := extract_from_map_by_value(fields[0], fields[1], data)
			if err != nil {
				return nil, errors.Wrap(err, "while extracting has_eq in dotpath case []interface{}")
			}
			return elems, nil
		}

		// handle listing
		if has_listing {
			elems, err := extract_from_map_by_field(inner, data)
			if err != nil {
				return nil, errors.Wrap(err, "while extracting listing in dotpath case []interface{}")
			}
			return elems, nil
		}

		// handle single string field/name
		val, ok := data[inner]
		if !ok {
			// try to look up by name
			name_value, ok := data["name"]
			if ok && name_value == inner {
				return data, nil
			}
			return nil, errors.New("could not find '" + inner + "' in object")
		}
		add_parent_and_path(val, data, path_str)
		if len(paths) == IDX+1 {
			return val, nil
		}
		ret, err := get_by_path(IDX+1, paths, val)
		if err != nil {
			return nil, errors.Wrapf(err, "from object "+inner)
		}
		return ret, nil
	}

	// handle default of field/name
	val, ok := data[P]
	if !ok {
		// try to look up by name
		name_value, ok := data["name"]
		if ok && name_value == P {
			return data, nil
		}
		return nil, errors.New("could not find '" + P + "' in object")
	}
	add_parent_and_path(val, data, path_str)
	if len(paths) == IDX+1 {
		return val, nil
	}
	ret, err := get_by_path(IDX+1, paths, val)
	if err != nil {
		return nil, errors.Wrapf(err, "from object "+P)
	}
	return ret, nil
}
