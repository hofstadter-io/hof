# Streaming Todo

This page is notes for when we stream `hof` development.


## Current Tasks

- [x] refactor lib/gen to use lib/runtime
- [ ] load datamodels & history during code gen      <====
  - [x] move logic for filtering gen>dm when calling `hof dm ...`
  - [ ] lift full history to DM if not already recorded there
  - [ ] funcs for handling the details of injecting into gen.In
  - [ ] call helper funcs at appropriate places (probably early)
- [ ] create a small test for gen+DM with history
- [ ] create database migrations using DM history
- [ ] getting-started/datamodel
- [ ] hofmod-sql & demo -> first-example(s)


Other

- [ ] Check out OCI code that was added to hof
