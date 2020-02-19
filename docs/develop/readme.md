# Development


### Building

`go build -mod=vendor`

Dependencies (vendor directory) are committed to git.

### Testing

Go to the examples directory and run the tool against one of the modules.

### Modifying the Grammar

See `pkg/parser/hof.peg`

##### Generating the Parser

[Pigeon](https://github.com/mna/pigeon):

```bash
# install pigeon
go get -u github.com/mna/pigeon

# generate parser
./ci/dev/gen-parser.sh
```


