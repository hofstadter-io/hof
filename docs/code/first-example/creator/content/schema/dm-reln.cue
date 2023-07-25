#Model: dm.#Model & {

	// Adds GoType
	Relations: [string]: R={
		GoType: string

		// Switch pattern
		GoType: [
			if R.Reln == "BelongsTo" {"*\(R.Type)"},
			if R.Reln == "HasOne" {"*\(R.Type)"},
			if R.Reln == "HasMany" {"[]*\(R.Type)"},
			if R.Reln == "ManyToMany" {"[]*\(R.Type)"},
			"unknown relation type: \(R.Reln)",
		][0]
	}
}
