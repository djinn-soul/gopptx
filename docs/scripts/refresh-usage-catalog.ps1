[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$buildScript = Join-Path $repoRoot "docs\scripts\build_usage_catalog.py"
$exportScript = Join-Path $repoRoot "scripts\tools\visual_regression\export_pptx_png.ps1"
$pptxDir = Join-Path $repoRoot "docs\assets\pptx\usage"
$pngDir = Join-Path $repoRoot "docs\assets\images\usage"
$tmpRoot = Join-Path $pngDir "tmp_catalog_export"

if (-not (Test-Path -LiteralPath $buildScript)) {
    throw "Missing build script: $buildScript"
}
if (-not (Test-Path -LiteralPath $exportScript)) {
    throw "Missing export script: $exportScript"
}

function Resolve-Python {
    $venvPython = Join-Path $repoRoot ".venv\Scripts\python.exe"
    if (Test-Path -LiteralPath $venvPython) { return @($venvPython) }

    $pythonCmd = Get-Command python -ErrorAction SilentlyContinue
    if ($pythonCmd) { return @("python") }

    $pyCmd = Get-Command py -ErrorAction SilentlyContinue
    if ($pyCmd) { return @("py", "-3") }

    throw "Python runtime not found. Install Python or use .venv/Scripts/python.exe."
}

$pythonExec = @(Resolve-Python)

if (@($pythonExec).Count -eq 1) {
    & $pythonExec[0] $buildScript
} else {
    & $pythonExec[0] $pythonExec[1] $buildScript
}

New-Item -ItemType Directory -Force -Path $pptxDir | Out-Null
New-Item -ItemType Directory -Force -Path $pngDir | Out-Null
New-Item -ItemType Directory -Force -Path $tmpRoot | Out-Null

$pptxFiles = Get-ChildItem -LiteralPath $pptxDir -Filter "*-python.pptx" | Sort-Object Name
if ($pptxFiles.Count -eq 0) {
    throw "No PPTX files found in $pptxDir"
}

# Remove stale non-catalog PNGs from older naming schemes.
Get-ChildItem -LiteralPath $pngDir -Filter "*.png" -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -notlike "*-python.png" } |
    Remove-Item -Force

foreach ($pptx in $pptxFiles) {
    $base = [IO.Path]::GetFileNameWithoutExtension($pptx.Name)
    $tmpOut = Join-Path $tmpRoot $base
    if (Test-Path -LiteralPath $tmpOut) {
        Remove-Item -LiteralPath $tmpOut -Recurse -Force
    }
    New-Item -ItemType Directory -Force -Path $tmpOut | Out-Null

    & $exportScript -PptxPath $pptx.FullName -OutDir $tmpOut

    $preferred = Get-ChildItem -LiteralPath $tmpOut -File -ErrorAction SilentlyContinue |
        Where-Object { $_.Name -match "^Slide2(\.|$)" } |
        Sort-Object Name |
        Select-Object -First 1
    if (-not $preferred) {
        $preferred = Get-ChildItem -LiteralPath $tmpOut -File -ErrorAction SilentlyContinue |
            Where-Object { $_.Name -match "^Slide1(\.|$)" } |
            Sort-Object Name |
            Select-Object -First 1
    }

    if (-not $preferred) {
        throw "No Slide1/Slide2 image found for $($pptx.Name)"
    }

    $destPng = Join-Path $pngDir ($base + ".png")
    Copy-Item -LiteralPath $preferred.FullName -Destination $destPng -Force
    Write-Output "Updated PNG: $destPng"
}

if (Test-Path -LiteralPath $tmpRoot) {
    Remove-Item -LiteralPath $tmpRoot -Recurse -Force
}

Write-Output "Usage catalog refresh complete."
Write-Output "PPTX count: $($pptxFiles.Count)"
Write-Output "PNG count:  $((Get-ChildItem -LiteralPath $pngDir -Filter '*.png').Count)"
