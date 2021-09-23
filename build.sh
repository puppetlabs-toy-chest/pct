#!/bin/bash

export WORKINGDIR=$(pwd)

target=${1:-build}
arch=$(go env GOHOSTARCH)
platform=$(go env GOHOSTOS)
binPath="$(pwd)/dist/pct_${platform}_${arch}"
binPath2="$(pwd)/dist/notel_pct_${platform}_${arch}"

if [ "$target" == "build" ]; then
  # Set goreleaser to build for current platform only
  if [ -z "${HONEYCOMB_API_KEY}" ]; then
    export HONEYCOMB_API_KEY="not_set"
  fi
  if [ -z "${HONEYCOMB_DATASET}" ]; then
    export HONEYCOMB_DATASET="not_set"
  fi
  goreleaser build --snapshot --rm-dist --single-target
  git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/baker-round "$binPath/templates"
  cp -r "$binPath/templates" "$binPath2/templates"
elif [ "$target" == "quick" ]; then
  go build -o ${binPath}/pct -tags telemetry
  go build -o ${binPath2}/pct
elif [ "$target" == "package" ]; then
  git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/baker-round "templates"
  goreleaser --skip-publish --snapshot --rm-dist
fi
