[CmdletBinding()]
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$pythonScript = Join-Path $repoRoot "docs\scripts\generate_s02_text_frame.py"
$pptxPath = Join-Path $repoRoot "docs\assets\pptx\usage\s02-text-frame-python.pptx"
$pngPath = Join-Path $repoRoot "docs\assets\images\usage\s02-text-frame-python.png"
$usageDir = Join-Path $repoRoot "docs\assets\images\usage"
$tmpDir = Join-Path $usageDir "tmp_s02_export"
$markdownPath = Join-Path $repoRoot "docs\showcase\usages\simple.md"
$exportScript = Join-Path $repoRoot "scripts\tools\visual_regression\export_pptx_png.ps1"

if (-not (Test-Path -LiteralPath $pythonScript)) {
    throw "Missing generator script: $pythonScript"
}
if (-not (Test-Path -LiteralPath $exportScript)) {
    throw "Missing export script: $exportScript"
}

function Resolve-Python {
    $venvPython = Join-Path $repoRoot ".venv\Scripts\python.exe"
    if (Test-Path -LiteralPath $venvPython) {
        return @($venvPython)
    }

    $pythonCmd = Get-Command python -ErrorAction SilentlyContinue
    if ($pythonCmd) {
        return @("python")
    }

    $pyCmd = Get-Command py -ErrorAction SilentlyContinue
    if ($pyCmd) {
        return @("py", "-3")
    }

    throw "Python runtime not found. Install Python or use .venv/Scripts/python.exe."
}

$pythonExec = @(Resolve-Python)

New-Item -ItemType Directory -Force -Path (Split-Path -Parent $pptxPath) | Out-Null
New-Item -ItemType Directory -Force -Path (Split-Path -Parent $pngPath) | Out-Null
New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null

if (@($pythonExec).Count -eq 1) {
    & $pythonExec[0] $pythonScript --out $pptxPath
} else {
    & $pythonExec[0] $pythonExec[1] $pythonScript --out $pptxPath
}

& $exportScript -PptxPath $pptxPath -OutDir $tmpDir

$slide1 = Get-ChildItem -LiteralPath $tmpDir -File -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -match "^Slide1(\.|$)" } |
    Sort-Object Name |
    Select-Object -First 1

if (-not $slide1) {
    # PowerPoint may export to parent folder on some hosts; fallback search.
    $slide1 = Get-ChildItem -LiteralPath $usageDir -File -ErrorAction SilentlyContinue |
        Where-Object { $_.Name -match "^Slide1(\.|$)" } |
        Sort-Object LastWriteTime -Descending |
        Select-Object -First 1
}

if (-not $slide1) {
    throw "Slide1 image not found in '$tmpDir' or '$usageDir'"
}

Copy-Item -LiteralPath $slide1.FullName -Destination $pngPath -Force

$s02Block = @'
## S02 - Basic Text Frame

**Focus:** Add controlled text regions and validate frame geometry with generated output.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
    p := pptx.NewPresentation()
    s := p.AddSlide()
    s.AddTextBox("Text Frame Properties")
    s.AddTextBox("0.5in margins demo")
    s.AddTextBox("Top anchor")
    s.AddTextBox("Bottom anchor")
    _ = p.Save("s02_text_frame_go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation

with Presentation.new("S02 Text Frame Demo") as p:
    p.add_slide("Text Frame Properties")
    p.add_shape(0, "rect", (40, 120, 180, 180), text="0.5in margins demo")
    p.add_shape(0, "rect", (240, 120, 180, 220), text="Top anchor")
    p.add_shape(0, "rect", (440, 120, 180, 220), text="Bottom anchor")
    p.add_shape(0, "rect", (40, 360, 180, 60), text="No wrap sample text")
    p.add_shape(0, "rect", (240, 360, 180, 100), text="Shrink-to-fit sample text")
    p.save("docs/assets/pptx/usage/s02-text-frame-python.pptx")
```

**Download PPTX:** [s02-text-frame-python.pptx](../../assets/pptx/usage/s02-text-frame-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Basic Text Frame](../../assets/images/usage/s02-text-frame-python.png)

'@

$md = Get-Content -Raw -Path $markdownPath
$pattern = '(?s)## S02 - Basic Text Frame.*?(?=## S03 - )'
$updated = [regex]::Replace($md, $pattern, $s02Block)
Set-Content -Path $markdownPath -Value $updated

if (Test-Path -LiteralPath $tmpDir) {
    Remove-Item -LiteralPath $tmpDir -Recurse -Force
}

# Cleanup stray PowerPoint exports if they landed in usage root.
Get-ChildItem -LiteralPath $usageDir -File -ErrorAction SilentlyContinue |
    Where-Object { $_.Name -match "^Slide\d+(\.|$)" } |
    Remove-Item -Force

Write-Output "Updated S02 artifacts and markdown:"
Write-Output "  PPTX: $pptxPath"
Write-Output "  PNG:  $pngPath"
Write-Output "  MD:   $markdownPath"
