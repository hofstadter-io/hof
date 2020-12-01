package langs

const CuelangModder = `
cue: {
	Name:        "cue"
	Version:     string | *"v0.3.0"
	ModFile:     string | * "cue.mods"
	SumFile:     string | * "cue.sums"
	ModsDir:     string | * "cue.mod/pkg"
	MappingFile: string | * "cue.mod/modules.txt"
	SymlinkLocalReplaces?: bool,
	InitTemplates: {...} | *{
		"cue.mod/module.cue": """
		module: "{{ .Module }}"

		"""
		...
	}
	VendorIncludeGlobs: [...string] | *[]
	VendorExcludeGlobs: [...string] | *[
		"/.git/**",
		"**/cue.mod/pkg/**",
	]
}
`
