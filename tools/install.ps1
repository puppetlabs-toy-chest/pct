function Install-Pct {
  [CmdletBinding()]
  param (
    [switch]$NoTelemetry
  )

  Set-StrictMode -Version 3.0
  $ErrorActionPreference = "Stop"

  $org = 'puppetlabs'
  $repo = 'pdkgo'

  $app = 'pct'

  $appPkgName = 'pct'
  if ($NoTelemetry) {
    $appPkgName = 'notel_pct'
  }

  $arch = "x86_64"
  $os = 'windows'
  $ext = '.zip'

  $ver = (Invoke-RestMethod "https://api.github.com/repos/${org}/${repo}/releases")[0].tag_name
  $file = "${appPkgName}_${ver}_${os}_${arch}${ext}"
  $downloadURL = "https://github.com/${org}/${repo}/releases/download/$ver/$file"

  $Destination = "~/.puppetlabs/pct"
  $Destination = $PSCmdlet.SessionState.Path.GetUnresolvedProviderPathFromPSPath($Destination)

  $tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ([System.IO.Path]::GetRandomFileName())
  $null = New-Item -ItemType Directory -Path $tempDir -Force -ErrorAction SilentlyContinue
  $packagePath = Join-Path -Path $tempDir -ChildPath $file

  if (!$PSVersionTable.ContainsKey('PSEdition') -or $PSVersionTable.PSEdition -eq "Desktop") {
    $oldProgressPreference = $ProgressPreference
    $ProgressPreference = "SilentlyContinue"
  }

  try {
    if ($NoTelemetry) {
      Write-Host "Downloading and extracting ${app} ${ver} (TELEMETRY DISABLED VERSION) to ${Destination}"
    } else {
      Write-Host "Downloading and extracting ${app} ${ver} to ${Destination}"
    }
    Invoke-WebRequest -Uri $downloadURL -OutFile $packagePath
  }
  finally {
    if (!$PSVersionTable.ContainsKey('PSEdition') -or $PSVersionTable.PSEdition -eq "Desktop") {
      $ProgressPreference = $oldProgressPreference
    }
  }

  if (Test-Path -Path $Destination) {
    Remove-Item -Path $Destination -Force -Recurse
  }
  Expand-Archive -Path $packagePath -DestinationPath $Destination

  Write-Host 'Remember to add the pct app to your path:'
  Write-Host "`$env:Path += `"`$env:PATH;${Destination}`""
}
