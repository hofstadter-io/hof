#!/usr/bin/env bash
set -euo pipefail

# make some users
curl -H 'Content-Type: application/json' -X POST -d '{ "Username": "bob", "Email": "bob@hof.io" }' localhost:8080/api/user
curl -H 'Content-Type: application/json' -X POST -d '{ "Username": "mary", "Email": "mary@hof.io" }' localhost:8080/api/user

# make some todos
curl -H 'Content-Type: application/json' -X POST -d '{ "Title": "hello", "Content": "world" }' localhost:8080/api/todo?username=bob
curl -H 'Content-Type: application/json' -X POST -d '{ "Title": "What did the fox say", "Content": "ringa ding ding something or other" }' localhost:8080/api/todo?username=bob

# add newline to printed output
echo
