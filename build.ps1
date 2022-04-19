#!/usr/bin/env pwsh

[CmdletBinding()]
param (
  [Parameter()]
  [ValidateSet('build', 'quick', 'package')]
  [string]
  $Target = 'build'
)
$Env:WORKINGDIR = $PSScriptRoot

$arch = go env GOHOSTARCH
$platform = go env GOHOSTOS
$binPath = Join-Path $PSScriptRoot "dist" "pct_${platform}_${arch}"
$binPath2 = Join-Path $PSScriptRoot "dist" "notel_pct_${platform}_${arch}"

$amd64 = go env GOAMD64
if ($amd64) {
	$binPath = "${binPath}_${amd64}"
	$binPath2 = "${binPath2}_${amd64}"
}

switch ($Target) {
  'build' {
    # Set goreleaser to build for current platform only
    # Add environment variables for honeycomb if not already loaded
    if (!(Test-Path ENV:\HONEYCOMB_API_KEY)) {
      $ENV:HONEYCOMB_API_KEY = 'not_set'
    }
    if (!(Test-Path ENV:\HONEYCOMB_DATASET)) {
      $ENV:HONEYCOMB_DATASET = 'not_set'
    }
    goreleaser build --snapshot --rm-dist --single-target
    git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/baker-round (Join-Path $binPath "templates")
		Get-ChildItem -Path (Join-Path $binPath "templates")
    Copy-Item (Join-Path $binPath "templates") -Destination (Join-Path $binPath2 "templates") -Recurse
  }
  'quick' {
    If ($Env:OS -match '^Windows') {
      go build -o "$binPath/pct.exe" -tags telemetry
      go build -o "$binPath2/pct.exe"
    } else {
      go build -o "$binPath/pct" -tags telemetry
      go build -o "$binPath2/pct"
    }
  }
  'package' {
    git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/baker-round "templates"
    goreleaser --skip-publish --snapshot --rm-dist
  }
}
