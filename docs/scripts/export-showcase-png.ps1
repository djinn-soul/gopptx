[CmdletBinding()]
param(
    [switch]$CleanTemp
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$exportScript = Join-Path $repoRoot "scripts\tools\visual_regression\export_pptx_png.ps1"

if (-not (Test-Path -LiteralPath $exportScript)) {
    throw "Export script not found: $exportScript"
}

$docsPptxDir = Join-Path $repoRoot "docs\assets\pptx"
$docsImgDir = Join-Path $repoRoot "docs\assets\images\showcase"
$tmpRoot = Join-Path $docsImgDir ".tmp_export"

if (-not (Test-Path -LiteralPath $docsPptxDir)) {
    throw "Docs PPTX directory not found: $docsPptxDir"
}
New-Item -ItemType Directory -Force -Path $docsImgDir | Out-Null
New-Item -ItemType Directory -Force -Path $tmpRoot | Out-Null

$mappings = @(
    @{ Pptx = "basic-generation.pptx"; OutPng = "basic-gen.png" },
    @{ Pptx = "basic_usage.pptx"; OutPng = "basic-usage-slide-1.png" },
    @{ Pptx = "rich-slide.pptx"; OutPng = "rich-slide.png" },
    @{ Pptx = "text-styling.pptx"; OutPng = "text-styling.png" },
    @{ Pptx = "chart-radar.pptx"; OutPng = "chart-radar.png" },
    @{ Pptx = "brand-reskin.pptx"; OutPng = "reskin-result.png" }
)

foreach ($item in $mappings) {
    $pptxPath = Join-Path $docsPptxDir $item.Pptx
    if (-not (Test-Path -LiteralPath $pptxPath)) {
        Write-Warning "Skipping missing PPTX: $pptxPath"
        continue
    }

    $tmpOut = Join-Path $tmpRoot ([IO.Path]::GetFileNameWithoutExtension($item.Pptx))
    if (Test-Path -LiteralPath $tmpOut) {
        Remove-Item -LiteralPath $tmpOut -Recurse -Force
    }
    New-Item -ItemType Directory -Force -Path $tmpOut | Out-Null

    & $exportScript -PptxPath $pptxPath -OutDir $tmpOut

    $slide = Get-ChildItem -LiteralPath $tmpOut -Filter "Slide1.*" -File |
        Sort-Object Name |
        Select-Object -First 1

    if (-not $slide) {
        throw "No Slide1 image found after export for: $pptxPath"
    }

    $dest = Join-Path $docsImgDir $item.OutPng
    Copy-Item -LiteralPath $slide.FullName -Destination $dest -Force
    Write-Output "Updated: $dest"
}

if ($CleanTemp -and (Test-Path -LiteralPath $tmpRoot)) {
    Remove-Item -LiteralPath $tmpRoot -Recurse -Force
}