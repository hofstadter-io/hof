
type DiffOp struct {
	Op        string
	Path      string
	Value     interface{}
	PrevValue interface{}
}

func typediff(original, current interface{}, basepath string) []DiffOp {
	var diffs []DiffOp

	// fmt.Println("typediff:", basepath)

	switch O := original.(type) {

	case map[string]interface{}:

		C, ok := current.(map[string]interface{})

		// if both are the same type, then go by key
		// else, just replace the value
		if ok {
			ds := diffMaps(O, C, basepath)
			diffs = append(diffs, ds...)
		} else {
			op := DiffOp{
				Op:        "update",
				Path:      basepath,
				Value:     C,
				PrevValue: O,
			}
			diffs = append(diffs, op)
		}

	case []interface{}:
		C, ok := current.([]interface{})

		// if both are the same type, then go by key
		// else, just replace the value
		if ok {
			ds := diffSlices(O, C, basepath)
			diffs = append(diffs, ds...)
		} else {
			op := DiffOp{
				Op:        "update",
				Path:      basepath,
				Value:     C,
				PrevValue: O,
			}
			diffs = append(diffs, op)
		}

	}

	return diffs
}

func diffMaps(orig, curr map[string]interface{}, basepath string) []DiffOp {
	var diffs []DiffOp

	// fmt.Println("diff'n maps", basepath)

	// Loop over Original
	for key, o := range orig {
		path := key
		if basepath != "" {
			path = basepath + "." + path
		}

		c, ok := curr[key]
		//  Remove if not-present
		if !ok {
			op := DiffOp{
				// Op: "remove",
				Op:    "delete",
				Path:  path,
				Value: o,
			}
			diffs = append(diffs, op)
		} else {

			// check if types are the same
			if reflect.TypeOf(o) == reflect.TypeOf(c) {
				// check type of 'c'
				switch c.(type) {
				case map[string]interface{}:
					ds := typediff(o.(map[string]interface{}), c, path)
					diffs = append(diffs, ds...)
				case []interface{}:
					ds := typediff(o.([]interface{}), c, path)
					diffs = append(diffs, ds...)
				default:
					if !reflect.DeepEqual(o, c) {
						op := DiffOp{
							Op:        "update",
							Path:      path,
							Value:     c,
							PrevValue: o,
						}
						diffs = append(diffs, op)
					}
				}
			} else {
				op := DiffOp{
					Op:        "update",
					Path:      path,
					Value:     c,
					PrevValue: o,
				}
				diffs = append(diffs, op)
			}

		}
	} // end loop over O

	// Loop over Current
	for key, c := range curr {
		path := key
		if basepath != "" {
			path = basepath + "." + path
		}

		_, ok := orig[key]
		// Add if not-present, otherwise ignore because we handled it above
		if !ok {
			op := DiffOp{
				Op:    "create",
				Path:  path,
				Value: c,
			}
			diffs = append(diffs, op)
		}

	}

	return diffs
}

func diffSlices(orig, curr []interface{}, basepath string) []DiffOp {
	var diffs []DiffOp

	// fmt.Println("diff'n []", basepath)

	// short-circuit (there's nothing)
	if len(orig) == 0 && len(curr) == 0 {
		return diffs
	}

	// short-circuit (all new)
	if len(orig) == 0 && len(curr) > 0 {
		for j, c := range curr {
			path := fmt.Sprintf("[%d]", j)
			if basepath != "" {
				path = basepath + "." + path
			}
			op := DiffOp{
				Op:    "create",
				Path:  path,
				Value: c,
			}
			diffs = append(diffs, op)
		}
		return diffs
	}

	// short-circuit (all gone)
	if len(orig) > 0 && len(curr) == 0 {
		for j, o := range orig {
			path := fmt.Sprintf("[%d]", j)
			if basepath != "" {
				path = basepath + "." + path
			}
			op := DiffOp{
				// Op: "remove",
				Op:    "delete",
				Path:  path,
				Value: o,
			}
			diffs = append(diffs, op)
		}
		return diffs
	}

	// short-circuit (new elem type)
	if reflect.TypeOf(orig[0]) != reflect.TypeOf(curr[0]) {
		op := DiffOp{
			Op:        "update",
			Path:      basepath,
			Value:     curr,
			PrevValue: orig,
		}
		diffs = append(diffs, op)
		return diffs
	}

	// figure out which slice diff func to call
	switch o := orig[0].(type) {
	case map[string]interface{}:
		_, ok := o["name"]
		if ok {
			ds := diffSlicesNamed(orig, curr, basepath)
			diffs = append(diffs, ds...)
		} else {
			ds := diffSlicesMaps(orig, curr, basepath)
			diffs = append(diffs, ds...)
		}

	default:
		ds := diffSlicesSimple(orig, curr, basepath)
		diffs = append(diffs, ds...)
	}

	return diffs
}

func diffSlicesSimple(orig, curr []interface{}, basepath string) []DiffOp {
	var diffs []DiffOp

	// fmt.Println("diff'n []string", basepath)

	// assume both are string elems

	// elems as maps
	oe, ce := map[string]bool{}, map[string]bool{}
	for _, o := range orig {
		oe[o.(string)] = true
	}
	for _, c := range curr {
		ce[c.(string)] = true
	}

	// fmt.Println("oe", oe)
	// fmt.Println("ce", ce)

	// calc missing
	rmv, add := map[string]bool{}, map[string]bool{}
	for key, _ := range oe {
		_, ok := ce[key]
		if !ok {
			rmv[key] = true
		}
	}
	for key, _ := range ce {
		_, ok := oe[key]
		if !ok {
			add[key] = true
		}
	}

	// fmt.Println("rmv", rmv)
	// fmt.Println("add", add)

	// do the diff
	oi, ci := 0, 0
	for oi < len(orig) && ci < len(curr) {
		ov, cv := orig[oi].(string), curr[ci].(string)

		// same element, so continue
		if ov == cv {
			oi += 1
			ci += 1
			continue
		}

		// calc a path off of the oi index
		path := fmt.Sprintf("[%d]", oi)
		if basepath != "" {
			path = basepath + "." + path
		}

		// remove or replace
		if ok := rmv[ov]; ok {
			// rmv & add
			if ok2 := add[cv]; ok2 {
				op := DiffOp{
					Op:        "update",
					Path:      path,
					Value:     cv,
					PrevValue: ov,
				}
				diffs = append(diffs, op)
				delete(rmv, ov)
				delete(add, cv)
				oi += 1
				ci += 1
				continue
			} else {
				// just remove
				op := DiffOp{
					// Op: "remove",
					Op:    "delete",
					Path:  path,
					Value: ov,
				}
				diffs = append(diffs, op)
				delete(rmv, ov)
				oi += 1
				continue
			}
		} else {
			// more of an insert

			op := DiffOp{
				// Op: "insert",
				Op:    "create",
				Path:  path,
				Value: cv,
			}
			diffs = append(diffs, op)
			delete(add, cv)
			ci += 1
			continue
		}

	}

	// remove any remaining elements
	for oi < len(orig) {
		ov := orig[oi].(string)
		path := fmt.Sprintf("[%d]", oi)
		if basepath != "" {
			path = basepath + "." + path
		}

		op := DiffOp{
			// Op: "remove",
			Op:    "delete",
			Path:  path,
			Value: ov,
		}
		diffs = append(diffs, op)
		delete(rmv, ov)
		oi += 1
	}

	for ci < len(curr) {
		cv := curr[ci].(string)
		path := fmt.Sprintf("[%d]", oi)
		if basepath != "" {
			path = basepath + "." + path
		}

		op := DiffOp{
			// Op: "append",
			Op:    "create",
			Path:  path,
			Value: cv,
		}
		diffs = append(diffs, op)
		delete(add, cv)
		ci += 1
	}

	// fmt.Println("DIFF: []maps", diffs)
	// fmt.Println("rmv", rmv)
	// fmt.Println("add", add)
	// fmt.Println("\n")

	return diffs
}

func diffSlicesNamed(orig, curr []interface{}, basepath string) []DiffOp {
	var diffs []DiffOp

	// fmt.Println("diff'n []maps", basepath)

	// elems as maps
	OS, CS := []map[string]interface{}{}, []map[string]interface{}{}
	for _, v := range orig {
		OS = append(OS, v.(map[string]interface{}))
	}
	for _, v := range curr {
		CS = append(CS, v.(map[string]interface{}))
	}

	// elems as maps
	oe, ce := map[string]bool{}, map[string]bool{}
	for _, o := range OS {
		oe[o["name"].(string)] = true
	}
	for _, c := range CS {
		ce[c["name"].(string)] = true
	}

	// fmt.Println("oe", oe)
	// fmt.Println("ce", ce)

	// calc missing
	rmv, add := map[string]bool{}, map[string]bool{}
	for key, _ := range oe {
		_, ok := ce[key]
		if !ok {
			rmv[key] = true
		}
	}
	for key, _ := range ce {
		_, ok := oe[key]
		if !ok {
			add[key] = true
		}
	}

	// fmt.Println("rmv", rmv)
	// fmt.Println("add", add)

	// do the diff
	oi, ci := 0, 0
	for oi < len(OS) && ci < len(CS) {
		ov, cv := OS[oi]["name"].(string), CS[ci]["name"].(string)

		// same element, so continue
		if ov == cv {
			// but first, recurse
			ds := typediff(OS[oi], CS[ci], basepath+"."+ov)
			diffs = append(diffs, ds...)
			oi += 1
			ci += 1
			continue
		}

		// calc a path off of the oi index
		path := fmt.Sprintf("[%d]", oi)
		if basepath != "" {
			path = basepath + "." + ov
		}

		// remove or replace
		if ok := rmv[ov]; ok {
			// rmv & add
			if ok2 := add[cv]; ok2 {
				op := DiffOp{
					Op:        "update",
					Path:      path,
					Value:     CS[ci],
					PrevValue: OS[oi],
				}
				diffs = append(diffs, op)
				delete(rmv, ov)
				delete(rmv, cv)
				oi += 1
				ci += 1
				continue
			} else {
				// just remove
				op := DiffOp{
					// Op: "remove",
					Op:    "delete",
					Path:  path,
					Value: OS[oi],
				}
				diffs = append(diffs, op)
				delete(rmv, ov)
				oi += 1
				continue
			}
		} else {
			// more of an insert

			path = basepath + "." + cv
			op := DiffOp{
				// Op: "insert",
				Op:    "create",
				Path:  path,
				Value: CS[ci],
			}
			diffs = append(diffs, op)
			delete(rmv, cv)
			ci += 1
			continue
		}

	}

	// remove any remaining elements
	for oi < len(OS) {
		ov := OS[oi]["name"].(string)
		path := ov
		if basepath != "" {
			path = basepath + "." + path
		}

		op := DiffOp{
			// Op: "remove",
			Op:    "delete",
			Path:  path,
			Value: OS[oi],
		}
		diffs = append(diffs, op)
		delete(rmv, ov)
		oi += 1
	}

	for ci < len(CS) {
		cv := CS[ci]["name"].(string)
		path := cv
		if basepath != "" {
			path = basepath + "." + path
		}

		op := DiffOp{
			// Op: "append",
			Op:    "create",
			Path:  path,
			Value: CS[ci],
		}
		diffs = append(diffs, op)
		delete(add, cv)
		ci += 1
	}

	// fmt.Println("DIFF: []maps", diffs)
	// fmt.Println("rmv", rmv)
	// fmt.Println("add", add)
	// fmt.Println("\n\n")
	return diffs
}

func diffSlicesMaps(orig, curr []interface{}, basepath string) []DiffOp {
	var diffs []DiffOp

	/*
		op := DiffOp {
			Op: "not-implemented",
			Path: basepath,
			Value: "diffSlicesMaps",
		}
		diffs = append(diffs, op)
	*/

	return diffs
}
