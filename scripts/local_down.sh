#!/usr/bin/env zsh

set -e pipefail

docker compose -f templates/docker-compose.yml down -v
