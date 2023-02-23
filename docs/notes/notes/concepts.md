

---

#### data model:

- abstract/conceptual, write in CUE, input to code gen
- Single Point of Truth (SPOT)
- checkpoint & track history (like db migrations)
- auto migration calculation & more with history
- it's like an intermediate representation
- data model + config / generator options

#### code gen:

- declarative programming (like IaC, DSLs)
- data + templates = _
- a generator is schemas & template
- any language or technology
- control over the generated code
- support integrated custom code
- regeneration when data model changes

User fills in data model & schema to gens used,

When combined with the data model history
even more complex code can be automated.

#### modules:

- any git repo with a `cue.mods` file
- any combo of data model & code gen
- all dep managed & composable
- framework to enable ecosystem

#### low/no code & dnd builder

- https://budibase.com/
- https://github.com/windmill-labs/windmill
