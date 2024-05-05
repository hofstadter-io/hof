---
title: Data Layer

brief: Define, validate, version, and migrate models across technologies

weight: 40
draft: true
---

{{<lead>}}
The data layer is a combination of schemas, config, and annotations
that are specially treated by `hof`. The two primary goals are

1. Provide consistent data models for downstream consumers
2. Enable history, diff, and migration features for version flexibility

The `hof datamodel` command and `schema/dm` are central to this, but are designed
in a way that allows you to customize and extend the built-in base.
{{</lead>}}

_Note, `hof dm` is shorthand for `hof datamodel`.

## Schemas

The core of `hof datamodel` is a set of schemas for adding metadate to a value.
These indicate the various node types that give structure to your datamodel.
The enables a flexible model that can still be used by the git-like
features for tracking history, showing diffs, and generating migration code.

There are also schemas for common datamodel formats (like SQL)
and enrichers for different languages (like Go & Python).

### Core Schema

These core schemas are metadata that `hof` recognizes and treats specially
to enable the `hof datamodel` commands.

<br>

{{<codePane file="code/hof-schemas/dm/dm.html" title="Datamodel Schemas">}}


### Common Formats

There are currently two core formats

1. {{<hof-gh-link path="schema/dm/fields/common.cue" >}} for common field types.
1. {{<hof-gh-link path="schema/dm/sql" >}} is the base for a relational datamodel.

{{<codePane file="code/hof-schemas/dm/fields/common.html" title="Common Field Schema">}}

### Enrichers

Enrichers extend or enhance a the schema to add language or library specifics.
For example:

- add language specific type used during code generation
- map field types to library specific types

Enrichers are the most common type of datamodel customization
as the target of code generation depends on your preferred tech stack.

See some examples here: {{<hof-gh-link path="schema/dm/enrichers">}}


## Commands and Example

This example will show you the basics of a datamodel
and the `hof dm` commands.

{{<codeInner lang="sh" title="hof dm -h (snippet)">}}
# Example Usage   (dm is short for datamodel)

  $ hof dm list   (print known data models)
  NAME         TYPE       VERSION  STATUS  ID
  Config       object     -        ok      Config
  MyDatamodel  datamodel  -        ok      datamodel-abc123

  $ hof dm tree   (print the structure of the datamodels)
  $ hof dm diff   (prints a tree based diff of the datamodel)
  $ hof dm checkpoint -m "a message about this checkpoint"
  $ hof dm log    (prints the log of changes from latest to oldest)

  You can also use the -d & -e flags to subselect datamodels and nested values

# Learn more:
  - https://docs.hofstadter.io/getting-started/data-layer/
  - https://docs.hofstadter.io/data-modeling/

Usage:
  hof datamodel [command]

Aliases:
  datamodel, dm

Available Commands:
  checkpoint  create a snapshot of the data model
  diff        show the current diff or between datamodel versions
  list        print available datamodels
  log         show the history for a datamodel
  tree        print datamodel structure as a tree
{{</codeInner>}}


### Create a Datamodel

We'll use a relational datamodel, typical of a database, for our example.

{{<codeInner lang="cue" title="datamodel.cue">}}
package datamodel

import (
  "github.com/hofstadter-io/hof/schema/dm/sql"
  "github.com/hofstadter-io/hof/schema/dm/fields"
)

// Traditional database model which maps onto tables & columns
Datamodel: sql.Datamodel & {
  // implied through definition, duplicated here for example clarity
  $hof: metadata: {
    id:   "datamodel-abc123"
    name: "MyDatamodel"
  }

  Models: {
    User: {
      Fields: {
        ID:        fields.UUID
        CreatedAt: fields.Datetime
        UpdatedAt: fields.Datetime
        DeletedAt: fields.Datetime

        email:    fields.Email
        username: fields.String
        password: fields.Password
        verified: fields.Bool
        active:   fields.Bool

        persona: fields.Enum & {
          Vals: ["guest", "user", "admin", "owner"]
          Default: "user"
        }
      }
    }
  }
}
{{</codeInner>}}

### Checkpoint a Datamodel


### Update a Datamodel


### View a Datamodel