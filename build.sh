#!/bin/bash

target=${1:-build}

clone_or_copy_templates() {
  if [ -z "$LOCAL_TEMPLATE_PATH" ]; then
    git clone -b main --single-branch https://github.com/puppetlabs/baker-round "$1"
  else
    test -d "$1" && rm -rf "$1"
    mkdir "$1"
    cp -R "$LOCAL_TEMPLATE_PATH/." "$1"
  fi
}

if [ "$target" == "build" ]; then
  arch=$(go env GOHOSTARCH)
  platform=$(go env GOHOSTOS)
  binPath="$(pwd)/dist/pct_${platform}_${arch}"
  # Set goreleaser to build for current platform only
  goreleaser build --snapshot --rm-dist --single-target
  clone_or_copy_templates "$binPath/templates"
elif [ "$target" == "package" ]; then
  clone_or_copy_templates "templates"
  goreleaser --skip-publish --snapshot --rm-dist
fi
