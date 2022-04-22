<#
.Description
A script to install the DevX documentation site in the docs\site directory and run it.
If the site is already installed, executing this script will run the site locally.
Changes and files added to the contents of the docs\md\content directory will be displayed on this
locally hosted site.
#>

param(
    [Alias('D')]
    [Switch]$BuildDrafts
)

# Check working directory is the root directory of the PCTepository
if (!(Test-Path -Path ".\docs\md")) {
    Throw 'Please run this script from the root directory of the PCTroject.'
}
# Check hugo extended is installed
if ([string]::IsNullOrEmpty((hugo version | Select-String -Pattern "extended"))) {
    Throw 'The "extended" version of hugo was not found, please install it.'
}

# Check that git and npm are installed
$Programs = @('git', 'npm')
foreach ($Program in $Programs) {
    try {
        Get-Command -Name $Program -ErrorAction Stop
    }
    catch {
        Throw "$Program was not found, please install it."
    }
}

# Check if the docs site directory exists
if (!(Test-Path -Path ".\docs\site")) {
    git clone https://github.com/puppetlabs/devx.git docs\site
    Push-Location docs\site
    Add-Content -Path .\go.mod -value 'replace github.com/puppetlabs/pctocs/md => ..\md'
    npm install
}
else {
    Push-Location docs\site
    git pull
    hugo mod clean
}

git submodule update --init --recursive --depth 50
hugo mod get
# Check if -D flag is set to true
try {
    if ($BuildDrafts) {
        hugo server -D
    }
    else {
        hugo server
    }
} catch {

} finally {
    Pop-Location
}
