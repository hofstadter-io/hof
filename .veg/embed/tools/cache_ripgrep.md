This tool uses ripgrep (`rg`) to match regular expression patterns in files.
Use this tool to quickly find identifiers, types, functions, defintions, etc... across many files.

The args have the following purpose

- `path` - base path to grep from
- `regexp` - regular expression to grep for
- `lines_around` - number of surrounding lines to include, max 8, defaults to 0 (only matching lines)
- `max_depth` - maximum recursive directory depth to grep, max 3, defaults to 1

The actual command run is `rg -Rn -d$max_depth -B$lines_around -A$lines_around --sort=path -e '$regexp' $path`