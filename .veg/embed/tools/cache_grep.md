Searches for content matching the given regular expression or literal string.
Uses Rust regex syntax; escape literal ., [, ], {, }, | with backslashes.
Use this tool to quickly find identifiers, types, functions, defintions, etc... across many files.

The args have the following purpose

- `path` - base path to grep from
- `glob` - glob pattern to match (e.g. "*.md")
- `regexp` - regular expression to grep for