# DMA - Your Data Model Assistant.

`dma` is a tool to help you manage your data
schemas, migrations, servers, and more.

Features:

- Schemas written in Cue. All the power of Cue and `dma` will turn into SQL,libraries, and tools.
- Migration calculation and management. `dma` tracks how your schema changes and handles the rest.
- Manage local, test, cicd, and production servers.
- Visual tool for interacting.
- Supports [list of databases]


## Command Structure

```
dma
- init
- config
- ui

- modelset
  - create
  - view
  - list
  - status
  - graph
  - diff
  - migrate
  - delete
  - test

- store
  - create
  - view
  - list
  - checkpoint
  - restore
  - delete

  # Things against live store (or run one)
  - run
  - status
  - conn
  - diff
  - up
  - down

- import
- export
- generate
```

## DMA directories and files

```
models/         (default) location for models, can have multiple
dma/            Dir for dba. You should commit this directory!
  meta.cue      model / store meta data
  modelsets/
    <modelset-name>/
      history/        # modelset history by subdirectory
      migrations/     # modelset migrations by subdirectory
  stores/
    <store-name>/
      data/           # for local running instances ???
```

## Cuelibs

- hofmod-cli
- hofmod-cuecfg
- hofmod-model
- structural

## Notes and references:

SQL Drivers / Libraries:

- https://pkg.go.dev/mod/github.com/lib/pq
- GORM
- https://github.com/jackc/pgx
- https://github.com/xo/dburl
- https://awesome-go.com/#database
- https://awesome-go.com/#database-drivers
- https://awesome-go.com/#orm

SQL Builders:

- https://github.com/go-jet/jet
- https://github.com/jirfag/go-queryset
- https://github.com/xo/xo

SQL Parsers:

- https://github.com/lfittl/pg_query_go
- https://gowalker.org/github.com/cockroachdb/cockroach/pkg/sql/parser
- https://marianogappa.github.io/software/2019/06/05/lets-build-a-sql-parser-in-go/
- https://github.com/pingcap/parser
- https://github.com/xwb1989/sqlparser

SQL -> Go:

- https://github.com/xo/xo
- https://github.com/kyleconroy/sqlc

Migration management:

- https://github.com/facebookincubator/ent
- https://www.prisma.io/docs/reference/tools-and-interfaces/prisma-migrate
- https://github.com/golang-migrate/migrate
- https://github.com/sequelize/umzug
- https://awesome-go.com/#database
- https://awesome-go.com/#orm

Other:

- https://github.com/vitessio/vitess
- https://github.com/google/go-cloud
- https://github.com/google/wire
- https://github.com/Azure/autorest


## Other ideas

Generate and/or manage seed data, database snapshots

- create snapshot from current state
- fill with snapshot

