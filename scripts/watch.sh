#!/usr/bin/env zsh

set -e pipefail

air --build.cmd "go build -o bin/orbis ./cmd/orbis" --build.bin "./bin/orbis"

