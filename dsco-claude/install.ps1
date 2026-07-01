#!/usr/bin/env pwsh
<#
.SYNOPSIS
  dsco-claude installer/updater for Windows PowerShell.

.DESCRIPTION
  Symlinks the dsco-expert agent (and any dsco-specific skills) into Claude
  Code. Falls back to copying when symlink creation is not permitted (no admin
  rights and Developer Mode off). On Linux/macOS/WSL/Git Bash use install.sh.

.PARAMETER Command
  install (default) | update | uninstall | status

.PARAMETER Project
  Target <dir>\.claude instead of ~\.claude. Pass "." for the current dir.

.PARAMETER Copy
  Copy files instead of symlinking.

.EXAMPLE
  ./install.ps1                 # install into ~\.claude
  ./install.ps1 status
  ./install.ps1 update
  ./install.ps1 -Project .      # install into .\.claude
  ./install.ps1 uninstall
#>
[CmdletBinding()]
param(
  [Parameter(Position = 0)]
  [ValidateSet('install', 'update', 'uninstall', 'status')]
  [string]$Command = 'install',
  [string]$Project,
  [switch]$Copy
)

$ErrorActionPreference = 'Stop'

$BundleDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$AgentSrc  = Join-Path $BundleDir 'agents\dsco-expert.md'
$SkillsSrc = Join-Path $BundleDir 'skills'
$Version   = (Get-Content (Join-Path $BundleDir 'VERSION') -ErrorAction SilentlyContinue |
              Select-Object -First 1)
if (-not $Version) { $Version = 'unknown' }

if ($PSBoundParameters.ContainsKey('Project')) {
  if ([string]::IsNullOrWhiteSpace($Project)) { $Project = '.' }
  $ClaudeDir = Join-Path (Resolve-Path $Project) '.claude'
} else {
  $ClaudeDir = Join-Path $HOME '.claude'
}

function Link-One($src, $dst) {
  $parent = Split-Path -Parent $dst
  if (-not (Test-Path $parent)) { New-Item -ItemType Directory -Force -Path $parent | Out-Null }
  if (Test-Path $dst) { Remove-Item -Recurse -Force $dst }
  if ($Copy) {
    Copy-Item -Recurse -Force $src $dst
    Write-Host "  copied   $dst"
  } else {
    try {
      New-Item -ItemType SymbolicLink -Path $dst -Target $src -Force | Out-Null
      Write-Host "  linked   $dst -> $src"
    } catch {
      Copy-Item -Recurse -Force $src $dst
      Write-Host "  symlink not permitted (enable Developer Mode or run as admin); copied instead: $dst"
    }
  }
}

function Each-Skill($action) {
  if (-not (Test-Path $SkillsSrc)) { return }
  Get-ChildItem -Directory $SkillsSrc | Where-Object {
    Test-Path (Join-Path $_.FullName 'SKILL.md')
  } | ForEach-Object { & $action $_.Name $_.FullName }
}

function Remove-One($path) {
  if (Test-Path $path) { Remove-Item -Recurse -Force $path; Write-Host "  removed  $path" }
}

function Status-One($path) {
  $item = Get-Item $path -ErrorAction SilentlyContinue
  if ($null -eq $item) {
    Write-Host "  missing  $path"
  } elseif ($item.LinkType -eq 'SymbolicLink') {
    Write-Host "  linked   $path -> $($item.Target)"
  } else {
    Write-Host "  present  $path (copy)"
  }
}

switch ($Command) {
  { $_ -in 'install', 'update' } {
    Write-Host "dsco-claude v$Version -> $ClaudeDir"
    if (Test-Path $AgentSrc) {
      Link-One $AgentSrc (Join-Path $ClaudeDir 'agents\dsco-expert.md')
    }
    Each-Skill { param($name, $src) Link-One $src (Join-Path $ClaudeDir "skills\$name") }
    Write-Host 'done.'
  }
  'uninstall' {
    Write-Host "dsco-claude: removing from $ClaudeDir"
    Remove-One (Join-Path $ClaudeDir 'agents\dsco-expert.md')
    Each-Skill { param($name, $src) Remove-One (Join-Path $ClaudeDir "skills\$name") }
    Write-Host 'done.'
  }
  'status' {
    Write-Host "dsco-claude v$Version"
    Write-Host "bundle:  $BundleDir"
    Write-Host "target:  $ClaudeDir"
    if (Test-Path $AgentSrc) { Status-One (Join-Path $ClaudeDir 'agents\dsco-expert.md') }
    Each-Skill { param($name, $src) Status-One (Join-Path $ClaudeDir "skills\$name") }
  }
}
