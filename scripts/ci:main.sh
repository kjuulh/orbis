#!/usr/bin/env zsh

set -e

dagger call --mod=./ci main --service="$SERVICE" --source=./ --version="${VERSION:-$(uuidgen)}"
