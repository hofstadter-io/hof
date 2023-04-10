---
title: Code

weight: 10
---


{{<lead>}}
This page provides an outline for developing and contributing to
`hof`'s code base.
{{</lead>}}

### Building


- `make hof` will go install the cli
- `hof gen` will regenerate code from the CUE in `design/`.
  You only have to run this if the design changes.
- `make formatters` will build local images.


### Testing

Most tests are Testscript base, a tool from the Go compiler
for testing CLIs. It makes it easy to write and add new tests.
The test runners live in `test.cue` and are run with `hof flow`.
These act as top-level wrappers around the individual tests,
and the same commands are run in CI for parity with local development.

- `hof flow @test/<area>` will run the tests for `<area>`
- `hof flow --list @test/*` will list available tests (areas)

You can find the paths from the `test.cue` file
to find where and how each set of tests are run.
Most tests are in the following directories:

- `test/...`
- `lib/mod/testdata`
- `lib/datamodel/test`
- `flow/testdata`
- `formatters/test`


### Core packages

- `cmd/hof/...` - the command structure, mostly generated
- `lib/hof` - metadata for values hof recognizes
- `lib/runtime` - holds CUE, generators, datamodels, workflows from entrypoints
- `lib/gen` - code generation engine
- `lib/template` - template rendering and helpers
- `lib/create` - creators implementation
- `lib/datamodel` - datamodel engine
- `formatters/tools` | `lib/fmt` - code and containers for code formatting
- `flow/...` - the workflow engine
- `lib/mod` | `lib/repos` - CUE dependency management
