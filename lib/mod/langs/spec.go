package langs

var ModderSpec = `
{
	[N=string]: {
		Name: N,
		Version: string,

		ModFile:  string,
		SumFile:  string,
		ModsDir:  string,
		MappingFile: string,

		NoLoad?: bool,
		CommandInit?: [...[...string]],
		CommandGraph?: [...[...string]],
		CommandTidy?: [...[...string]],
		CommandVendor?: [...[...string]],
		CommandVerify?: [...[...string]],
		CommandStatus?: [...[...string]],

		InitTemplates?: {
			[string]: string
		},
		InitPreCommands?: [...[...string]],
		InitPostCommands?: [...[...string]],

		VendorIncludeGlobs?: [...string],
		VendorExcludeGlobs?: [...string],
		VendorTemplates?: {
			[string]: string
		},
		VendorPreCommands?: [...[...string]],
		VendorPostCommands?: [...[...string]],

		ManageFileOnly?: bool,
		SymlinkLocalReplaces?: bool,

		IntrospectIncludeGlobs?: [...string],
		IntrospectExcludeGlobs?: [...string],
		IntrospectExtractRegex?: [...string],

		PackageManagerDefaultPrefix?: string,
	}
}
`
