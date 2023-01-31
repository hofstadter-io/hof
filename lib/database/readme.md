# Data Systems utilities

Used in `datastore`, `datamodel` & `flow`


## Connectors

### SQL

Goal is to be CGO free

- PostgreSQL: https://github.com/jackc/pgx
- SQLite: https://gitlab.com/cznic/sqlite
- MySQL: https://github.com/go-mysql-org/go-mysql
- MSSQL: https://github.com/denisenkom/go-mssqldb
- Oracle: https://github.com/sijms/go-ora
- Snowflake: https://github.com/snowflakedb/gosnowflake

Other

- rest / introspection: github.com/go-rest/rest
- migrations: https://github.com/golang-migrate/migrate
- psql migs: https://github.com/jackc/tern

## NoSQL

Mongo, Elastic


## Object storage

S3, GCS, ...

(probably just in `flow/tasks/objs`)


## Cache

Redis, memcache
(probably just in `flow/tasks/kv`)

## Messaging

Watermill (https://github.com/ThreeDotsLabs/watermill/) covers most options

(probably just in `flow/tasks/msg`)

