---
title: Upgrading from 0.6.7 to 0.6.8
weight: 11
---

{{<lead>}}
This page will guide you on upgrading through breaking changes.
If you find any that are missing, please reach out on GitHub or Slack.
{{</lead>}}

## Upgrading your Datamodels

<br>

The schemas for `hof datamodel` and `github.com/hofstadter-io/hof/schema/dm/...` both changed
more than we can provide backwards compatibility for.

The primary goals were:

- to remove our structure and patterns from the core data model
- to enable arbtirary structure defined by the user
- to enable arbitrary and namespaced enrichments
- to make template authoring easier

The primary changes were:

- minimize the core [data model schema](https://github.com/hofstadter-io/hof/blob/_dev/schema/dm/dm.cue#L28)
- move the previous structure to a [schema/dm/sql package](https://github.com/hofstadter-io/hof/blob/_dev/schema/dm/sql/dm.cue)
- change how relations between models is represented


### Upgrading Imports

We changed two things that impact imports

1. reorganizing the schema as just described
2. changing the schemas from CUE definitions to stucts for openness

Before, you had a single import:

```go
import "github.com/hofstadter-io/hof/schemas/dm"

#MyModels: dm.#Datamodel & { ... }
```

After, you will likely need multiple imports:

```go
import (
    "github.com/hofstadter-io/hof/schemas/dm"
    "github.com/hofstadter-io/hof/schemas/dm/sql"
    "github.com/hofstadter-io/hof/schemas/dm/fields"
)

#MyModels: sql.Datamodel & { ... }
```

Note the removal of the `#` in `sql.Datamodel`.
The same happened for generators but we were able to provide backwards compatibility there.
You should still remove the `#` from any `gen.#Generator` so they become `gen.Generator`.


### Upgrading Relations

The way you specify relations has changed.

1. There is no longer a `Model.Relations` field to hold them all. Instead, the `Model.Fields.[name].Relation` is used.
1. The name of the fields that hold relation details have changed.
1. The casing of the relation type has change from TitleCase to kebab-case

The primary reason for this change was to make template authoring easier.
You now only need one loop to process all model fields,
rather than before where you needed a loop each for `Fields` and `Relations`.
You can still get each separately by filtering the `Fields` on the `Relation` sub-field.

<br>

#### Before ([old schema](https://github.com/hofstadter-io/hof/blob/v0.6.7/schema/dm/dm.cue))

```text
#MyModels: dm.#Datamodel & {

    Models: {
        User: {
            Fields: { 
                Username: { Type: "string" }
            }
            Relations: {
                Posts: {
                   Reln: "HasMany"
                   Type: "Post"
                }
            }
        }
        Post: {
            Fields: {
                Title: { Type: "string" }
                Body:  { Type: "string" }
            }
            Relations: {
                Author: {
                   Reln: "BelongsTo"
                   Type: "User"
                }
            }
        }
    }
}
```

<br>

#### Changes

- `dm.#Datamodel` -> `sql.Datamodel`
- `Model.Relations.[name].Reln` -> `Model.Fields.[name].Relation.Type`
- `Model.Relations.[name].Type` -> `Model.Fields.[name].Relation.Other`

<br>

#### After ([new schema](https://github.com/hofstadter-io/hof/blob/v0.6.8/schema/dm/sql/dm.cue))

```text
MyModels: sql.Datamodel & {

    Models: {
        User: {
            Fields: { 
                Username: { Type: "string" }
                Posts: {
                    Type: "string"
                    Relation: {
                        Reln: "has-many"
                        Type: "Post"
                    }
                }
            }
        }
        Post: {
            Fields: {
                Title:  { Type: "string" }
                Body:   { Type: "string" }
                Author: {
                    Type: "string"
                    Relation: {
                        Reln: "has-many"
                        Type: "Post"
                    }
                }
            }
        }
    }
}
```

