package dotpath

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	log "gopkg.in/inconshreveable/log15.v2" // logging framework
)

var logger = log.New()

func init() {
	config_logger("warn")
}

func SetLogger(l log.Logger) {
	lcfg := viper.GetStringMap("log-config.dotpath.default")

	if lcfg == nil || len(lcfg) == 0 {
		logger = l
	} else {
		level_str := lcfg["level"].(string)
		stack := lcfg["stack"].(bool)
		level, err := log.LvlFromString(level_str)
		if err != nil {
			panic(err)
		}

		termlog := log.LvlFilterHandler(level, log.StdoutHandler)
		if stack {
			term_stack := log.CallerStackHandler("%+v", log.StdoutHandler)
			termlog = log.LvlFilterHandler(level, term_stack)
		}

		logger.SetHandler(termlog)
	}
}

func SetLogLevel(level string) {
	config_logger(level)
}

func config_logger(level string) {
	logger = log.New()
	term_level, err := log.LvlFromString(level)
	if err != nil {
		panic(err)
	}

	term_stack := log.CallerStackHandler("%+v", log.StdoutHandler)
	term_caller := log.CallerFuncHandler(log.CallerFileHandler(term_stack))
	termlog := log.LvlFilterHandler(term_level, term_caller)

	/*
		term_caller := log.CallerFuncHandler(log.CallerFileHandler(log.StdoutHandler))
		termlog := log.LvlFilterHandler(term_level, term_caller)
	*/

	//	termlog := log.LvlFilterHandler(term_level, log.StdoutHandler)
	logger.SetHandler(termlog)

}

func Get(path string, data interface{}, no_solo_array bool) (interface{}, error) {
	if data == nil {
		return nil, errors.New("input data is nil")
	}
	if path == "." {
		return data, nil
	}
	paths := strings.Split(path, ".")
	if len(paths) < 1 {
		return nil, errors.New("Bad path supplied: " + path)
	}
	if strings.Contains(paths[0], ":") {
		pos := strings.Index(paths[0], ":")
		paths[0] = paths[0][pos+1:]
	}

	// fmt.Println("GETPATH:", path, paths, data)

	ret, err := get_by_path(0, paths, data)
	if err != nil {
		return nil, err
	}
	if T, ok := ret.([]interface{}); ok {
		if no_solo_array && len(T) == 1 {
			return T[0], nil
		}
	}
	return ret, nil
}

func GetByPathSlice(path []string, data interface{}, no_solo_array bool) (interface{}, error) {
	if data == nil {
		return nil, errors.New("input data is nil")
	}
	if path[0] == "." {
		return data, nil
	}
	if strings.Contains(path[0], ":") {
		pos := strings.Index(path[0], ":")
		path[0] = path[0][pos+1:]
	}
	ret, err := get_by_path(0, path, data)
	if err != nil {
		return nil, err
	}
	if T, ok := ret.([]interface{}); ok {
		if no_solo_array && len(T) == 1 {
			return T[0], nil
		}
	}
	return ret, nil
}

func get_by_path(IDX int, paths []string, data interface{}) (interface{}, error) {
	if data == nil {
		return nil, errors.New("input data is nil")
	}
	header := fmt.Sprintf("get_by_path:  %d  %v  in:\n%+v\n\n", IDX, paths, data)
	// fmt.Println(header)
	logger.Debug(header)

	P := paths[IDX]

	lpos_index := strings.Index(P, "[")
	rpos_index := strings.LastIndex(P, "]")
	pos_colon := strings.Index(P, ":")
	has_listing := strings.Contains(P, ",")
	pos_regex := strings.Index(P, "regex")

	has_eq := strings.Contains(P, "==")
	has_ne := strings.Contains(P, "!=")
	has_ge := strings.Contains(P, ">=")
	has_gt := strings.Contains(P, ">")
	has_le := strings.Contains(P, "<=")
	has_lt := strings.Contains(P, "<")

	logger.Debug("Has: ", "idx", IDX, "curr", P, "paths", paths,
		"lpos", lpos_index, "rpos", rpos_index, "slicing", pos_colon,
		"listing", has_listing, "regex", pos_regex,
	)
	logger.Debug("has bool",
		"has_eq", has_eq, "has_ne", has_ne,
		"has_ge", has_ge, "has_gt", has_gt,
		"has_le", has_le, "has_lt", has_lt,
	)

	switch T := data.(type) {

	case map[string]interface{}:
		elems, err := get_from_smap_by_path(IDX, paths, T)
		if err != nil {
			return nil, errors.Wrap(err, "while extracting path from smap in get_by_path")
		}
		// if E, ok := elems.([]interface{}); ok && len(E) == 1 {
		// 	return E[0], nil
		// }
		return elems, nil

	case map[interface{}]interface{}:
		elems, err := get_from_imap_by_path(IDX, paths, T)
		if err != nil {
			return nil, errors.Wrap(err, "while extracting path from imap in get_by_path")
		}
		// if E, ok := elems.([]interface{}); ok && len(E) == 1 {
		// 	return E[0], nil
		// }
		return elems, nil

	case []interface{}:
		logger.Debug("Processing Slice", "paths", paths, "T", T)
		elems, err := get_from_slice_by_path(IDX, paths, T)
		if err != nil {
			return nil, errors.Wrap(err, "while extracting path from slice")
		}
		if len(paths) == IDX+1 {
			return elems, nil
		} else {
			switch E := elems.(type) {
			case []interface{}:
				ees := []interface{}{}
				for _, e := range E {
					ee, eerr := get_by_path(IDX+1, paths, e)
					if eerr == nil {
						ees = append(ees, ee)
					}
				}
				if len(ees) == 1 {
					return ees[0], nil
				}
				return ees, nil
			default:
				ees, eerr := get_by_path(IDX+1, paths, elems)
				if eerr != nil {
					return nil, errors.Wrap(eerr, "while extracting path from slice in default")
				}
				//if E, ok := ees.([]interface{}); ok && len(E) == 1 {
				//	return E[0], nil
				//}
				return ees, nil
			}
		}

	default:
		typ := reflect.TypeOf(data)
		if typ.Kind() == reflect.Ptr {
			d := reflect.ValueOf(data)
			data = reflect.Indirect(d).Interface()
			typ = reflect.TypeOf(data)
			logger.Debug("pointer dereference", "data", data, "type", typ, "d", d)
		}

		switch typ.Kind() {
		case reflect.Map:
			M := reflect.ValueOf(data)
			logger.Debug("Map: ", "map", M, "type", typ)
			v := M.MapIndex(reflect.ValueOf(P))
			if IDX+1 >= len(paths) {
				return v.Interface(), nil
			}

			ees, eerr := get_by_path(IDX+1, paths, v.Interface())
			if eerr != nil {
				return nil, errors.Wrap(eerr, "while extracting path from slice in default")
			}
			// if E, ok := ees.([]interface{}); ok && len(E) == 1 {
			// 	return E[0], nil
			// }
			return ees, nil

		case reflect.Slice:
			S := reflect.ValueOf(data)
			logger.Debug("Slice: ", "slice", S, "type", typ)
			if IDX+1 >= len(paths) {
				return S.Interface(), nil
			}
			ees, eerr := get_by_path(IDX+1, paths, S.Interface())
			if eerr != nil {
				return nil, errors.Wrap(eerr, "while extracting path from slice in default")
			}
			// if E, ok := ees.([]interface{}); ok && len(E) == 1 {
			// 	return E[0], nil
			// }
			return ees, nil

		case reflect.Struct:

			S := reflect.ValueOf(data)
			logger.Debug("Struct: ", "struct", S, "type", typ)
			F := S.FieldByName(P)
			if IDX+1 >= len(paths) {
				return F.Interface(), nil
			}

			ees, eerr := get_by_path(IDX+1, paths, F.Interface())
			if eerr != nil {
				return nil, errors.Wrap(eerr, "while extracting path from slice in default")
			}
			// if E, ok := ees.([]interface{}); ok && len(E) == 1 {
			// 	return E[0], nil
			// }
			return ees, nil

		default:
			str := fmt.Sprintf("%+v", typ)
			return nil, errors.New("unknown data object type: " + str)

		}
		str := fmt.Sprintf("%+v", typ)
		return nil, errors.New("unknown data object type: " + str)

	} // END of type switch

}

func add_parent_and_path(child interface{}, parent interface{}, path string) (interface{}, error) {
	logger.Debug("adding parent to child", "child", child, "parent", parent, "path", path)
	parent_ref := "unknown-parent"
	switch P := parent.(type) {

	case map[string]interface{}:
		p_ref, ok := P["name"]
		if !ok {
			return nil, errors.Errorf("parent does not have name: %+v", parent)
		}
		parent_ref = p_ref.(string)
	case map[interface{}]interface{}:
		p_ref, ok := P["name"]
		if !ok {
			return nil, errors.Errorf("parent does not have name: %+v", parent)
		}
		parent_ref = p_ref.(string)

	default:
		str := fmt.Sprintf("%+v", reflect.TypeOf(parent))
		return nil, errors.New("unknown parent object type: " + str)

	}

	switch C := child.(type) {

	case map[string]interface{}:
		C["parent"] = parent_ref
		C["path"] = path
	case map[interface{}]interface{}:
		C["parent"] = parent_ref
		C["path"] = path

	case []interface{}:
		for _, elem := range C {
			switch E := elem.(type) {
			case map[string]interface{}:
				E["parent"] = parent_ref
				E["path"] = path
			case map[interface{}]interface{}:
				E["parent"] = parent_ref
				E["path"] = path
			default:
				str := fmt.Sprintf("in slice of %+v", reflect.TypeOf(E))
				return nil, errors.New("element not an object type: " + str)
			}
		}

	default:
		str := fmt.Sprintf("%+v", reflect.TypeOf(C))
		return nil, errors.New("unknown data object type: " + str)

	}
	return child, nil
}
