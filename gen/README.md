# This directory contains generated code

This gen directory exists because concepts in `hof`,
which overlap with capabilities we want to generate,
would collide if we generated into the lib directory.
So we create a directory specificly as:

- `lib` is all hand written
- `gen` is largely generated
- `cmd` is also largely generated
