package dotpath

import (
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
)

func extract_from_slice_with_splice(splice string, data []interface{}) (interface{}, error) {

	// handle slicing
	fields := strings.Split(splice, ":")

	// get slicing values
	l, r := -1, -1
	var ierr error
	if f := fields[0]; f != "" {
		l, ierr = strconv.Atoi(f)
		if ierr != nil {
			return nil, errors.Wrapf(ierr, "converting lpos in path in extract_splice: "+splice)
		}
	}
	if f := fields[1]; f != "" {
		r, ierr = strconv.Atoi(f)
		if ierr != nil {
			return nil, errors.Wrapf(ierr, "converting rpos in path in extract_splice: "+splice)
		}
	}

	// fmt.Println("L,R: ", l, r, data)

	// do things based on positions
	if l > -1 && r > -1 {
		return data[l:r], nil
	} else if l == -1 && r > -1 {
		return data[:r], nil
	} else if l > -1 && r == -1 {
		return data[l:], nil
	}
	return data, nil
}

func extract_from_slice_with_name(listing string, data []interface{}) (interface{}, error) {

	ret := []interface{}{}
	// separate listing
	fields := strings.Split(listing, ",")

	sort.Strings(fields)

	for _, d := range data {
		name := ""
		switch D := d.(type) {
		case map[string]interface{}:
			name_value, ok := D["name"].(string)
			if ok {
				name = name_value
			}

		case map[interface{}]interface{}:
			name_value, ok := D["name"].(string)
			if ok {
				name = name_value
			}

		}

		if name != "" {
			pos := sort.SearchStrings(fields, name)
			if pos < len(fields) && fields[pos] == name {
				ret = append(ret, d)
			}
		}
	}
	return ret, nil
}

func extract_from_slice_with_field(field, value string, data []interface{}) (interface{}, error) {

	ret := []interface{}{}
	// separate listing
	fields := strings.Split(value, ",")
	sort.Strings(fields)

	for _, d := range data {
		field_value := ""
		switch D := d.(type) {
		case map[string]interface{}:
			f_value, ok := D[field].(string)
			if ok {
				field_value = f_value
			}

		case map[interface{}]interface{}:
			f_value, ok := D[field].(string)
			if ok {
				field_value = f_value
			}

		}

		if field_value != "" {
			pos := sort.SearchStrings(fields, field_value)
			if pos < len(fields) && fields[pos] == field_value {
				ret = append(ret, d)
			}
		}
	}
	return ret, nil
}

func extract_listing_from_map_string(listing string, data map[string]interface{}) (interface{}, error) {

	ret := []interface{}{}
	// separate listing
	fields := strings.Split(listing, ",")

	sort.Strings(fields)

	for _, f := range fields {
		val, ok := data[f]
		if ok {
			ret = append(ret, val)
		}
	}
	return ret, nil
}

func extract_listing_from_map_iface(listing string, data map[interface{}]interface{}) (interface{}, error) {

	ret := []interface{}{}
	// separate listing
	fields := strings.Split(listing, ",")

	sort.Strings(fields)

	for _, f := range fields {
		val, ok := data[f]
		if ok {
			ret = append(ret, val)
		}
	}
	return ret, nil
}

func extract_from_map_by_field(field string, data interface{}) (interface{}, error) {

	ret := []interface{}{}
	// separate listing
	fields := strings.Split(field, ",")

	sort.Strings(fields)

	switch D := data.(type) {
	case map[string]interface{}:
		for _, f := range fields {
			val, ok := D[f]
			if ok {
				ret = append(ret, val)
			}
		}

	case map[interface{}]interface{}:
		for _, f := range fields {
			val, ok := D[f]
			if ok {
				ret = append(ret, val)
			}
		}

	default:
		return nil, errors.New("data arg is not a map type")
	}

	if len(ret) == 1 {
		return ret[0], nil
	}
	return ret, nil
}

func extract_from_map_by_value(field interface{}, value string, data interface{}) (interface{}, error) {

	ret := []interface{}{}
	// separate listing
	fields := strings.Split(value, ",")

	sort.Strings(fields)

	switch D := data.(type) {
	case map[string]interface{}:
		for _, f := range fields {
			val, ok := D[f]
			if ok && val == value {
				ret = append(ret, val)
			}
		}

	case map[interface{}]interface{}:
		for _, f := range fields {
			val, ok := D[f]
			if ok && val == value {
				ret = append(ret, val)
			}
		}

	default:
		return nil, errors.New("data arg is not a map type")
	}

	fmt.Println("RETTTTTT:", ret)

	if len(ret) == 1 {
		return ret[0], nil
	}
	return ret, nil
}
