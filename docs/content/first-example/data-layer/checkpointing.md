---
title: Checkpointing
brief: and Data Model history
weight: 60
---

{{<lead>}}
This section is largely a light introduction to the `hof datamodel` command, or `hof dm` for short.
{{</lead>}}

As we have mentioned, the `hof dm` subcommands are utilities for managing the versioning and history of your datamodels.

{{<codeInner lang="sh">}}
Available Commands:
  checkpoint  create a snapshot of the data model
  diff        show the diff between data model version
  history     list the snapshots for a data model
  info        print details for a data model
  list        print available data models
  log         show the history of diffs for a data model
{{</codeInner>}}

Run `hof dm -h` for the full listing.
The top-level [data-modeling](/reference/hof-datamodel/) goes deeper into
the command and many related topics.

### Checkpointing

We can run `hof dm checkpoint` to create
a versioned snapshot of our data model.
With many snapshots, we can introspect
the history and diff between versions.
This opens up some interesting possibilities:

- check backwards compatibility and suggest appropriate semantic version changes
- automatically generating migrations as our program and data model evolve
- generating transformation functions between versions
    - for supporting many versions while keeping business logic clear of these details
		- feature flag capabilities
		- deploying migrations and applications updates without downtime or coordination

The inspiration for some of this is

- automating database migrations as applications evolve
- [Project Cambria](https://www.inkandswitch.com/cambria/) for the lenses concept
- [grafana/thema](https://github.com/grafana/thema) for lacuna, or filling the gaps


### Trying `hof dm`

You can try out the `hof dm` commands by

1. `hof dm checkpoint` to create a first snapshot
2. `hof dm ...` list, info, history to introspect
3. add some fields or new models
4. `hof dm diff` to see what has changed
5. the introspection commands to see backwards compatibility


