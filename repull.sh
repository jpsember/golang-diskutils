#!/usr/bin/env bash
set -u

rm go.mod
cp _SKIP_go.mod go.mod
go build
go mod tidy
