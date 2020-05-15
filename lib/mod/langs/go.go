package langs

const GolangModder = `
go: {
	Name:          "go",
	Version:       string | *"1.14",
	ModFile:       string | *"go.mod",
	SumFile:       string | *"go.sum",
	ModsDir:       string | *"vendor",
	MappingFile:   string | *"vendor/modules.txt",
	CommandInit:   [...[...string]] | *[["go", "mod", "init"]],
	CommandGraph:  [...[...string]] | *[["go", "mod", "graph"]],
	CommandTidy:   [...[...string]] | *[["go", "mod", "tidy"]],
	CommandVendor: [...[...string]] | *[["go", "mod", "vendor"]],
	CommandVerify: [...[...string]] | *[["go", "mod", "verify"]],
}
`
