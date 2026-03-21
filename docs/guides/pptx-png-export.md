# Docs Screenshot Pipeline (PPTX -> PNG)

This page is for documentation maintainers.
It explains how docs screenshots are generated from PPTX files for pages in `docs/showcase/`.

## Purpose

- keep screenshots in docs aligned with real generated presentations
- avoid stale or manually edited images
- make visual updates reproducible in one command

## Source and Target Paths

- PPTX sources: `docs/assets/pptx/`
- PNG outputs: `docs/assets/images/showcase/`
- Export helper: `docs/scripts/export-showcase-png.ps1`

## One-Command Refresh

```powershell
./docs/scripts/export-showcase-png.ps1 -CleanTemp
```

## What Gets Updated

The script maps these PPTX files to docs images:

- `basic-generation.pptx` -> `basic-gen.png`
- `basic_usage.pptx` -> `basic-usage-slide-1.png`
- `rich-slide.pptx` -> `rich-slide.png`
- `text-styling.pptx` -> `text-styling.png`
- `chart-radar.pptx` -> `chart-radar.png`
- `brand-reskin.pptx` -> `reskin-result.png`

## Requirements

- Windows host with Microsoft PowerPoint installed (COM automation used)
- existing `.pptx` artifacts under `docs/assets/pptx/`

## When To Run

- after updating example code and regenerating PPTX docs artifacts
- before opening docs PRs that modify screenshots
- when visual regressions are reported in showcase pages