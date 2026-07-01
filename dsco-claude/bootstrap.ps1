#!/usr/bin/env pwsh
<#
.SYNOPSIS
  dsco-claude bootstrap for Windows PowerShell: download the bundle from GitHub
  and install the skills into Claude Code, without a manual checkout.

.DESCRIPTION
  One-liner:
    irm https://raw.githubusercontent.com/byte4ever/dsco/master/dsco-claude/bootstrap.ps1 | iex

  To pin a version or pass options, run it as a scriptblock:
    & ([scriptblock]::Create((irm <url>/bootstrap.ps1))) -Ref v1.4.0-rc.1 -Copy

  Uses tar.exe and curl/Invoke-WebRequest, which ship with Windows 10+.

.PARAMETER Ref
  git ref to fetch (branch / tag / sha). Default: master (or $env:DSCO_CLAUDE_REF).

.PARAMETER Dir
  Where to place the bundle. Default: ~\.dsco-claude (or $env:DSCO_CLAUDE_HOME).

.PARAMETER Copy
  Copy files instead of symlinking (passed through to install.ps1).
#>
[CmdletBinding()]
param(
  [string]$Ref = $(if ($env:DSCO_CLAUDE_REF) { $env:DSCO_CLAUDE_REF } else { 'master' }),
  [string]$Dir = $(if ($env:DSCO_CLAUDE_HOME) { $env:DSCO_CLAUDE_HOME } else { Join-Path $HOME '.dsco-claude' }),
  [switch]$Copy
)

$ErrorActionPreference = 'Stop'
$repo = 'byte4ever/dsco'
$url  = "https://github.com/$repo/archive/$Ref.tar.gz"
$tmp  = Join-Path ([System.IO.Path]::GetTempPath()) ("dsco-claude-" + [System.Guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Force -Path $tmp | Out-Null

try {
  Write-Host "dsco-claude: fetching $repo@$Ref ..."
  $tar = Join-Path $tmp 'bundle.tar.gz'
  Invoke-WebRequest -UseBasicParsing -Uri $url -OutFile $tar
  tar -xzf $tar -C $tmp

  $src = Get-ChildItem -Path $tmp -Directory |
    ForEach-Object { Join-Path $_.FullName 'dsco-claude' } |
    Where-Object { Test-Path $_ } |
    Select-Object -First 1
  if (-not $src) {
    throw "dsco-claude/ not found in $repo@$Ref (the bundle ships on master and on releases that include it)."
  }

  if (-not (Test-Path $Dir)) { New-Item -ItemType Directory -Force -Path $Dir | Out-Null }
  foreach ($item in 'skills', 'install.sh', 'install.ps1', 'bootstrap.sh', 'bootstrap.ps1', 'VERSION', 'README.md', 'CHANGELOG.md') {
    $p = Join-Path $Dir $item
    if (Test-Path $p) { Remove-Item -Recurse -Force $p }
  }
  Copy-Item -Recurse -Force (Join-Path $src '*') $Dir

  Write-Host "dsco-claude: bundle placed in $Dir"
  $ps1 = Join-Path $Dir 'install.ps1'
  if ($Copy) { & $ps1 install -Copy } else { & $ps1 install }
} finally {
  Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
