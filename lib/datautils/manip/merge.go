package manip

import (
	"fmt"

	"github.com/pkg/errors"
)

/*
Where's your docs doc?!
*/
func Merge(original interface{}, update interface{}) (merged interface{}, err error) {

	if original == nil {
		return update, nil
	}

	if update == nil {
		return original, nil
	}

	// call the recursive merge
	return merge(original, update)
}

/*
Where's your docs doc?!
*/
func merge(original interface{}, update interface{}) (merged interface{}, err error) {

	logger.Info("Merging", "original", original, "update", update)

	switch O := original.(type) {

	case map[string]interface{}:
		U, ok := update.(map[string]interface{})
		if !ok {
			return nil, errors.New("update is not mS like original")
		}
		logger.Info("mS entering")
		for key, val := range U {
			logger.Debug("in merge mS-U", "key", key, "val", val, "curr", O[key])
			if curr, exists := O[key]; exists {
				tmp, err := merge(curr, val)
				logger.Debug("after merge mS", "tmp", tmp, "err", err)
				if err != nil {
					return nil, errors.Wrap(err, "in merge mS")
				}
				O[key] = tmp
			} else {
				O[key] = val
			}
		}
		logger.Info("mS returning", "O", O)
		return O, nil

	case []interface{}:
		U, ok := update.([]interface{})
		if !ok {
			return nil, errors.New("update is not aI like original")
		}
		// logger.Warn("O", "data", O)
		// logger.Warn("U", "data", U)

		logger.Info("aI entering")
		// turn update into map
		UM := map[string]interface{}{}
		for i, elem := range U {
			switch E := elem.(type) {

			case map[string]interface{}:
				name, ok := E["name"]
				if !ok {
					return nil, errors.New("original array objects must have names to be merged")
				}
				UM[name.(string)] = E

			case string:
				UM[E] = E

			default:
				logger.Error("original unknown elem type in aI", "i", i, "elem", elem)
				return nil, errors.New("original unknown elem type in aI")
			}
		}

		for i, elem := range O {
			// logger.Crit("O-loop", "i", i, "elem", elem)
			switch E := elem.(type) {

			case map[string]interface{}:
				iname, ok := E["name"]
				if !ok {
					return nil, errors.New("original array objects must have names to be merged")
				}

				name := iname.(string)
				// logger.Error("Name", "name", name)

				curr, exists := UM[name]
				if exists {
					tmp, err := merge(elem, curr)
					// this is correct, the var names curr and elem are backwards...
					// busy fixing a bug
					// logger.Crit("merging with existing element", "key", name, "val", curr, "curr", elem)
					if err != nil {
						return nil, errors.Wrap(err, "in merge MS")
					}
					O[i] = tmp
					delete(UM, name)
				}
			case string:
				_, exists := UM[E]
				if exists {
					delete(UM, E)
				}

			default:
				logger.Error("original unknown elem type in aI", "i", i, "elem", elem)
				return nil, errors.New("original unknown elem type in aI")
			}
		}
		// merge
		logger.Info("aI")

		// turn back into array
		OA := []interface{}{}
		for _, val := range O {
			OA = append(OA, val)
		}
		for _, elem := range U {
			switch E := elem.(type) {

			case map[string]interface{}:
				name, ok := E["name"]
				if !ok {
					return nil, errors.New("original array objects must have names to be merged")
				}
				_, exists := UM[name.(string)]
				if exists {
					OA = append(OA, elem)
				}

			case string:
				_, exists := UM[E]
				if exists {
					OA = append(OA, elem)
				}

			}
		}

		// logger.Error("OA", "data", OA)

		logger.Info("aI returning", "OA", OA)
		return OA, nil

	case string:
		return update, nil

	default:
		return nil, errors.New("unmergable original" + fmt.Sprintf("%t, %+v", original, original))

	}

	logger.Crit("Shouldn't get here (end of merge function)")
	return nil, errors.New("PANIC, should not get here")
}
