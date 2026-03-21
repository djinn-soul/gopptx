# Go API Reference

This page is the detailed reference for the Go surface of `gopptx`.
It is organized around the package-level entry points and the fluent builders that most Go users touch first.

Primary source files:

- `pkg/pptx/presentation.go`
- `pkg/pptx/slide.go`
- `pkg/pptx/presentation_builder.go`
- `pkg/pptx/templates/templates.go`
- `pkg/pptx/table_compat.go`
- `pkg/pptx/chart_compat.go`
- `pkg/pptx/connector.go`
- `pkg/pptx/elements/slide.go`

## Quick Start

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
    slides := []pptx.SlideContent{
        pptx.NewSlide("Intro").AddBullet("Hello from gopptx"),
    }

    data, err := pptx.CreateWithSlides("Deck Title", slides)
    if err != nil {
        panic(err)
    }

    _ = data
}
```

## Package-level constructors

These are the package-level entry points in `pkg/pptx`.

The detailed pages are split by concern:

- [Go Metadata, Protection, and Export Reference](go-metadata-export.md)
- [Go Slides Reference](go-slides.md)
- [Go Notes, Comments, and Sections Reference](go-notes-comments-sections.md)
- [Go Tables Reference](go-tables.md)
- [Go Templates Reference](go-templates.md)
- [Go Charts Reference](go-charts.md)
- [Go Shapes and Connectors Reference](go-shapes.md)

## PresentationBuilder

The builder is the fluent Go authoring path for full decks.

### `NewPresentationBuilder(title string) *PresentationBuilder`

Create a builder with a deck title.

### Builder methods

- `WithMetadata(meta Metadata) *PresentationBuilder`
- `WithSlideSize(size SlideSize) *PresentationBuilder`
- `WithTheme(theme styling.Theme) *PresentationBuilder`
- `WithMaster(master *elements.SlideMaster) *PresentationBuilder`
- `AddSlide(slide SlideContent) *PresentationBuilder`
- `AddTitleSlide(title string) *PresentationBuilder`
- `AddBulletSlide(title string, bullets []string) *PresentationBuilder`
- `AddShapesSlide(title string, shapes ...Shape) *PresentationBuilder`
- `AddCustomXML(content string) *PresentationBuilder`
- `WithSlideNumbers(show bool) *PresentationBuilder`
- `WithFooter(text string) *PresentationBuilder`
- `WithDateTime(show bool) *PresentationBuilder`
- `WithModifyPassword(password string) *PresentationBuilder`
- `WithMarkAsFinal(final bool) *PresentationBuilder`
- `WithSignaturesEnabled(enabled bool) *PresentationBuilder`
- `WithEncryptionPassword(password string) *PresentationBuilder`
- `Build() ([]byte, error)`
- `WriteToFile(path string) error`

## See also

- [API Reference](../api-reference.md)
- [Python Presentation API](python-presentation-api.md)
- [JSON Bridge Operations](bridge-operations.md)
