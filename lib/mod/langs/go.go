package langs

const GolangModder = `
langs: go: {
	Name:          "go",
	Version:       "1.16",
	ModFile:       "go.mod",
	SumFile:       "go.sum",
	ModsDir:       "vendor",
	MappingFile:   "vendor/modules.txt",
	CommandInit:   [["go", "mod", "init"]],
	CommandGraph:  [["go", "mod", "graph"]],
	CommandTidy:   [["go", "mod", "tidy"]],
	CommandVendor: [["go", "mod", "vendor"]],
	CommandVerify: [["go", "mod", "verify"]],
}
`
