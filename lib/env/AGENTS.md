# lib/env

`lib/env` implements a CUE interface to Dagger via their Go SDKs.

important directories:

- `schemas/env` - CUE definitions for all the things
- `lib/env/dag` - core that make CUE -> Dagger
- `lib/env/cmd` - bridge between `lib/env/dag` and main cli cobra commands
- `lib/env/devex` - working example for development
- `lib/env/incept` - hack to run hof with itself for pretty logs

## Layout

## Types and Schema


Naming Patterns:

- `<key>` -> `step<Key>`
  - `step<Key>Config`
- `#<key>` -> `hash<Key>`
  - `hash<Key>Config`
  - `hash<Key>Index`