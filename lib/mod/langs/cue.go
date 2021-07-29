package langs

const CuelangModder = `
langs: cue: {
	Name:        "cue"
	Version:     "v0.4.0"
	ModFile:     "cue.mods"
	SumFile:     "cue.sums"
	ModsDir:     "cue.mod/pkg"
	MappingFile: "cue.mod/modules.txt"
	InitTemplates: {
		"cue.mod/module.cue": """
		module: "{{ .Module }}"
		"""
	}
	VendorExcludeGlobs: [
		"/.git/**",
		"**/cue.mod/pkg/**",
	]
}
`
