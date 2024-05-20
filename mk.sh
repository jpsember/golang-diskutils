#!/usr/bin/env bash
set -eu

# dgen.sh
go build -o diskutils
go test ./...

