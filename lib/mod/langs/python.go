package langs

const PythonModder = `
python: {
	Name:          "python",
	Version:       string | *"3.8",
	ModFile:       string | *"python.mod",
	SumFile:       string | *"requirements.txt",
	ModsDir:       string | *"vendor",
	MappingFile:   string | *"vendor/modules.txt",

	CommandInit:   [...[...string]] | *[["python", "-m", "venv", "venv"]],
	CommandVendor: [...[...string]] | *[["bash", "-c", ". ./venv/bin/activate && pip install -r requirements.txt"]],
}
`
