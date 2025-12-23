#!/usr/bin/env bash
set -euo pipefail


# Top-level list
(
    echo "# dagger.io/dagger"
    echo ""
    echo '```go'
    go doc dagger.io/dagger 
    echo '```'
    echo ""

) > dag-gdoc-all.md

# Functions for core types
(
    items=(
        "Client"
        "Container"
        "Directory"
        "File"
        "Service"
    )

    for I in ${items[@]}; do
        echo "# dagger.$I"
        echo ""
        echo '```go'
        go doc dagger.io/dagger.$I
        echo '```'
        echo ""
        echo ""
    done
) > dag-gdoc-core-types.md