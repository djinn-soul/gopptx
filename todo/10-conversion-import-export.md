# 10 - Conversion / Import / Export

Source: `TODO.md` plus `web_sid/docs.aspose.com/slides/python-net/{convert-presentation,convert-slide,import-presentation,render-presentation-with-fallback-font,render-a-slide-as-an-svg-image,extract-text-from-presentation}`.

## Supported

| Item | Docs | Status | Notes |
| --- | --- | --- | --- |
| Export to PDF | `convert-presentation/` | Supported | PDF export is marked `[x]`. |
| Export to HTML5 | `convert-presentation/` | Supported | HTML5 export is marked `[x]`. |
| Render slide thumbnails | `convert-slide/`, `presentation-viewer/` | Supported | Thumbnail-style render helpers are marked `[x]`. |
| Render slide as SVG | `render-a-slide-as-an-svg-image/` | Supported | SVG rendering is marked `[x]`. |
| Render with fallback fonts | `render-presentation-with-fallback-font/` | Supported | Fallback-aware rendering/export is marked `[x]`. |

## Not Supported

| Item | Docs | Status | Notes |
| --- | --- | --- | --- |
| Export to image | `convert-slide/`, `convert-presentation/` | Not supported | PNG/JPEG export remains `[ ]`. |
| Export to video | `convert-presentation/` | Not supported | MP4 export remains `[ ]`. |
| Export to XAML | `convert-presentation/` | Not supported | XAML export remains `[ ]`. |
| Convert to / save as ODP | `convert-presentation/`, `convert-ppt-to-pptx/` | Not supported | ODP convert/save flows remain `[ ]`. |
| Import from ODP / PPT / external sources | `import-presentation/` | Not supported | ODP, legacy PPT, Google Slides, and Keynote import remain `[ ]`. |
| Convert PPT <-> PPTX | `ppt-vs-pptx/` | Not supported | PPT-to-PPTX and PPTX-to-PPT remain `[ ]`. |
| Export to Markdown / XPS / SWF / TIFF / JPG / PNG | `convert-presentation/`, `convert-slide/` | Not supported | These legacy export targets remain `[ ]`. |
| Dedicated slide-range export APIs | `convert-slide/` | Not supported | Single-slide/range export helpers remain `[ ]`. |
