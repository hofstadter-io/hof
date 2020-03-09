#!/usr/bin/env sh

go build cmd/lexer/main.go && mv main lexer
go build cmd/parser/main.go && mv main parser
