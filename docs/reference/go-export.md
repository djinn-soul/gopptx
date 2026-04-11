# Go Export Reference

This page documents the PDF and HTML export functions in `pkg/pptx/export`.

Primary source files:

- `pkg/pptx/export/pdf.go`
- `pkg/pptx/export/pdf_options.go`
- `pkg/pptx/export/html.go`
- `pkg/pptx/export/html_options.go`

> **Note:** The export package is imported separately as `github.com/djinn-soul/gopptx/pkg/pptx/export`.
> It is not re-exported from `pkg/pptx`.

## PDF export

### Functions

- `PDF(title string, slides []elements.SlideContent, outputPath string) error`
- `PDFWithOptions(title string, slides []elements.SlideContent, outputPath string, opts PDFOptions) error`
- `PDFFromFile(pptxPath, pdfPath string) error` — convert a saved PPTX file to PDF
- `PDFFromFileWithOptions(pptxPath, pdfPath string, opts PDFOptions) error`

### `PDFOptions`

```go
type PDFOptions struct {
    Driver          PDFDriver
    NativeFontPaths []string  // additional font search paths for the native driver
}
```

### `PDFDriver` constants

| Constant | Value | Description |
|---|---|---|
| `PDFDriverAuto` | `"auto"` | Tries LibreOffice/PowerPoint first, falls back to native |
| `PDFDriverNative` | `"native"` | Built-in renderer (experimental — lower fidelity for complex decks) |
| `PDFDriverLibreOffice` | `"libreoffice"` | Requires LibreOffice on `$PATH` |
| `PDFDriverPowerPoint` | `"powerpoint"` | Requires PowerPoint on macOS or Windows |

### `ParsePDFDriver(value string) (PDFDriver, error)`

Validates and normalizes a driver name string. Useful when reading driver from config/env.

### Typical usage

```go
import (
    "github.com/djinn-soul/gopptx/pkg/pptx"
    "github.com/djinn-soul/gopptx/pkg/pptx/export"
)

slides := []pptx.SlideContent{
    pptx.NewSlide("Title").AddBullet("First point"),
}

// Auto driver (default)
if err := export.PDF("My Deck", slides, "out.pdf"); err != nil {
    panic(err)
}

// Force native renderer with a custom font path
opts := export.PDFOptions{
    Driver:          export.PDFDriverNative,
    NativeFontPaths: []string{"/usr/share/fonts/custom"},
}
if err := export.PDFWithOptions("My Deck", slides, "out.pdf", opts); err != nil {
    panic(err)
}

// Convert an existing PPTX
if err := export.PDFFromFile("existing.pptx", "out.pdf"); err != nil {
    panic(err)
}
```

---

## HTML export

### Functions

- `HTML(title string, slides []elements.SlideContent) string` — returns HTML as a string
- `HTMLWithOptions(title string, slides []elements.SlideContent, opts HTMLOptions) string`
- `HTMLToWriter(w io.Writer, title string, slides []elements.SlideContent, opts HTMLOptions) error` — stream to any `io.Writer`
- `DefaultHTMLOptions() HTMLOptions`

### `HTMLOptions`

```go
type HTMLOptions struct {
    EmbedImages         bool         // inline images as base64 data URIs (default: true)
    BaseURL             string       // path prefix for sidecar image assets
    Theme               *ThemeColors // override CSS color variables
    IncludeNavigation   bool         // add prev/next slide buttons (default: true)
    IncludeSlideNumbers bool         // show slide number overlay (default: true)
}
```

### `ThemeColors`

```go
type ThemeColors struct {
    TitleColor      string
    BodyColor       string
    AccentColor     string
    BackgroundColor string
    SlideBackground string
}
```

### Typical usage

```go
import (
    "os"
    "github.com/djinn-soul/gopptx/pkg/pptx"
    "github.com/djinn-soul/gopptx/pkg/pptx/export"
)

slides := []pptx.SlideContent{
    pptx.NewSlide("My Deck").AddBullet("Point one"),
}

html := export.HTML("My Deck", slides)
_ = os.WriteFile("deck.html", []byte(html), 0o600)

// Or stream directly
f, _ := os.Create("deck.html")
defer f.Close()
_ = export.HTMLToWriter(f, "My Deck", slides, export.DefaultHTMLOptions())
```

## See also

- [Go API Reference](go-api.md)
- [Go Slides Reference](go-slides.md)
