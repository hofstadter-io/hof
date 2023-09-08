package browser

func (V *Browser) Encode() (map[string]any, error) {
	m := map[string]any{
		"type": V.TypeName(),
		"mode": V.mode,
		"usingScope": V.usingScope,
		"docs": V.docs,
		"attrs": V.attrs,
		"defs": V.defs,
		"optional": V.optional,
		"ignore": V.ignore,
		"inline": V.inline,
		"resolve": V.resolve,
		"concrete": V.concrete,
		"hidden": V.hidden,
		"final": V.final,
		"validate": V.validate,
	}

	var err error
	m["source"], err = V.source.Encode()
	return m, err
}


