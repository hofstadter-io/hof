#!/bin/bash
set -euo pipefail

echo ""
echo 'DOWNLOAD Test dependencies - go get em killer!'
echo "===================================="
go get -u github.com/golang/lint/golint                        # Linter
go get -u honnef.co/go/tools/cmd/megacheck                     # Badass static analyzer/linter
go get -u github.com/fzipp/gocyclo
go get -u github.com/franela/goblin                            # Better testing framework

