---
title: Data Layer

brief: Define, validate, version, and migrate models across technologies

weight: 40
---

{{<lead>}}
The data layer is a combination of schemas, config, and annotations
that `hof` understands and operates on. The primary goals are:

- Support git-style checkpointing, history, and diff features, flexible to your datamodel schema
- Provide consistent data models for downstream consumers
- Enable downstream features like automatic database migrations, client/server version skew, and change detection for infrastructure as code

The `hof datamodel` command and `schema/dm` schema form the foundations and 
are designed so you can customize, extend, or replace as needed.

- the built-in base models, fields, and enrichers
- the shape and hierarchy for diff and history tracking

{{</lead>}}

_Note, `hof dm` is shorthand for `hof datamodel`_.

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

To create a datamodel, simply write some CUE. 
To see your datamodel in hof, run the following commands

```text
hof dm list
hof eval datamodel.cue
hof tui view datamodel.cue
```

<details>
<summary>
hof dm list
</summary>
{{<codeInner>}}
NAME       TYPE       VERSION  STATUS      ID        
Datamodel  datamodel  -        no-history  Datamodel
{{</codeInner>}}
</details>

<details>
<summary>
hof eval datamodel.cue
</summary>
{{<codePane file="code/getting-started/data-layer/create/hof-eval.html" >}}
</details>
<br>

{{<codePane title="datamodel.cue" file="code/getting-started/data-layer/create/datamodel.html" >}}


### Checkpoints and History

Like a database and SQL migration files, you can checkpoint the history of your datamodels.
This is an optional feature, but will allow you to automatically generate database migrations
and code that can upgrade requests or downgrade responses, allowing for client/server version skew.

To checkpoint a datamodel, run `hof dm checkpoint -s ... -m "..."`


{{<codeInner>}}
> hof dm checkpoint --suffix initial_user_model --message "initial user model for the application"
creating checkpoint: 20240507010604_initial_user_model "initial user model for the application"
{{</codeInner>}}

At the root of your CUE module, you should now find a `.hof/dm/...` directory

{{<codeInner>}}
> tree .hof 
.hof
└── dm
    └── Datamodel
        ├── 20240507010604_initial_user_model.cue
        └── Models
            └── User
                └── 20240507010604_initial_user_model.cue

5 directories, 2 files
{{</codeInner>}}

The `hof` SQL datamodel tracks both the full datamodel and the individual models.
This is done to ease the authoring of code generation templates that create
database migrations and version skew functions.
Generally, `hof` supports user defined datamodel hierarchy and history tracking.


### Update a Datamodel

Next, we will add a UserProfile to the model.

<details>
<summary>
hof eval datamodel.cue
</summary>
{{<codePane file="code/getting-started/data-layer/update/hof-eval.html" >}}
</details>
<br>

{{<codePane title="datamodel.cue" file="code/getting-started/data-layer/update/datamodel.html" >}}

### View a Datamodel

With our modified datamodel, we can explore some `hof dm` commands for inspecting it.

We can see that it has changes with `hof dm list`, note the `dirty` status.

{{<codeInner title="hof dm list">}}
NAME       TYPE       VERSION  STATUS  ID        
Datamodel  datamodel  -        dirty   Datamodel
{{</codeInner>}}

We can also see the diff of those changes.
`hof` uses a structural diff on the CUE value
which allows for the hierarchical history.

<details>
<summary>
hof dm diff
</summary>
{{<codePane file="code/getting-started/data-layer/update/hof-diff.html" >}}
</details>
<br>

With a new checkpoint...

{{<codeInner title="hof dm checkpoint...">}}
> hof dm checkpoint -s add_user_profile -m "add a user profile and give ownership to the user"
creating checkpoint: 20240507014051 "add a user profile and give ownership to the user"
{{</codeInner>}}

we can also view the history log

{{<codeInner>}}
> hof dm log
20240507014051_add_user_profile: "add a user profile and give ownership to the user"
  Datamodel         ~ has changes
    Models
      User          ~ has changes
      UserProfile   + new value

20240507010604_initial_user_model: "initial user model"
  Datamodel         + new value
    Models
      User          + new value
{{</codeInner>}}


If we inspect the `.hof/dm` directory, we will see there are three new files.
One for the datamodel change, and one for each model that was changed.

{{<codeInner>}}
> tree .hof
.hof
└── dm
    └── Datamodel
        ├── 20240507010604_initial_user_model.cue
        ├── 20240507014051_add_user_profile.cue
        └── Models
            ├── User
            │   ├── 20240507010604_initial_user_model.cue
            │   └── 20240507014051_add_user_profile.cue
            └── UserProfile
                └── 20240507014051_add_user_profile.cue

6 directories, 5 files
{{</codeInner>}}
