# Go API Reference

This page is the detailed reference for the Go surface of `gopptx`.
It is organized around the package-level entry points and the fluent builders that most Go users touch first.

Primary source files:

- `pkg/pptx/presentation.go`
- `pkg/pptx/presentation_api.go`
- `pkg/pptx/presentation_chart_api.go`
- `pkg/pptx/presentation_chart_data_builder.go`
- `pkg/pptx/presentation_builder.go`
- `pkg/pptx/slide.go`
- `pkg/pptx/elements/slide.go`
- `pkg/pptx/slide_animation.go`
- `pkg/pptx/hyperlink.go`
- `pkg/pptx/image.go`
- `pkg/pptx/placeholder.go`
- `pkg/pptx/layout_helpers.go`
- `pkg/pptx/markdown_compat.go`
- `pkg/pptx/table_compat.go`
- `pkg/pptx/chart_compat.go`
- `pkg/pptx/connector.go`
- `pkg/pptx/templates/templates.go`
- `pkg/pptx/styling/` (colors, fonts, units, themes)
- `pkg/pptx/export/` (PDF and HTML — imported separately)

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

- [Go Metadata, Protection, and Export Reference](go-metadata-export.md) — `Open*`, `Presentation` runtime, core properties, images, markdown, placeholders
- [Go Slides Reference](go-slides.md) — `SlideContent`, transitions, animations, backgrounds, masters, SmartArt, layout constants
- [Go Notes, Comments, and Sections Reference](go-notes-comments-sections.md) — builder notes + `PresentationEditor` comments/sections
- [Go Tables Reference](go-tables.md) — `Table`, `TableCell` including all fluent methods
- [Go Templates Reference](go-templates.md)
- [Go Charts Reference](go-charts.md) — constructors, `Series` type, data builders, runtime chart API
- [Go Shapes and Connectors Reference](go-shapes.md) — all shape constructors including callouts, flowchart extended set, badge, images
- [Go Text Reference](go-text.md) — `Run`, `Paragraph`, `ParagraphStyle` and all rich-text constants
- [Go Hyperlinks Reference](go-hyperlinks.md) — `NewHyperlink`, all `HyperlinkAction` constructors
- [Go Export Reference](go-export.md) — PDF and HTML export, `PDFOptions`, `HTMLOptions`
- [Go Layout Helpers Reference](go-layout-helpers.md) — `Center`, `Grid`, `Stack`, `DistributeUniform`

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

## Styling: Length units

Source file: `pkg/pptx/styling/units.go`

- `Inches(value float64) Length`
- `Centimeters(value float64) Length`
- `Points(value float64) Length`
- `Emu(value int64) Length`
- `FontSize(pt float64) int`
- `(l Length).Inches() float64`
- `(l Length).Cm() float64`
- `(l Length).Pt() float64`
- `(l Length).Emu() int64`

### Font size presets

`FontSizeTitle` (44), `FontSizeSubtitle` (32), `FontSizeHeading` (28), `FontSizeBody` (18), `FontSizeSmall` (14), `FontSizeCaption` (12), `FontSizeCode` (14), `FontSizeLarge` (36), `FontSizeXLarge` (48)

## Styling: Colors

Source file: `pkg/pptx/styling/colors.go`

Basic: `ColorRed`, `ColorGreen`, `ColorBlue`, `ColorWhite`, `ColorBlack`, `ColorGray`, `ColorLightGray`, `ColorDarkGray`, `ColorYellow`, `ColorLightBlue`, `ColorOrange`, `ColorPurple`, `ColorCyan`, `ColorMagenta`, `ColorNavy`, `ColorTeal`, `ColorOlive`

Corporate: `ColorCorporateBlue`, `ColorCorporateGreen`, `ColorCorporateRed`, `ColorCorporateOrange`

Material Design: `ColorMaterialRed`, `ColorMaterialPink`, `ColorMaterialPurple`, `ColorMaterialIndigo`, `ColorMaterialBlue`, `ColorMaterialCyan`, `ColorMaterialTeal`, `ColorMaterialGreen`, `ColorMaterialLime`, `ColorMaterialAmber`, `ColorMaterialOrange`, `ColorMaterialBrown`, `ColorMaterialGray`

IBM Carbon: `ColorCarbonBlue60`, `ColorCarbonBlue40`, `ColorCarbonGray100`, `ColorCarbonGray80`, `ColorCarbonGray20`, `ColorCarbonGreen50`, `ColorCarbonRed60`, `ColorCarbonPurple60`

## Styling: Line dash constants

Source file: `pkg/pptx/styling/line_style.go`

`LineDashSolid`, `LineDashDash`, `LineDashDot`, `LineDashDashDot`, `LineDashDashDotDot`, `LineDashLongDash`, `LineDashLongDashDot`

## Styling: Built-in themes

Source file: `pkg/pptx/styling/theme.go`

### Theme variables

`ThemeCorporate`, `ThemeModern`, `ThemeVibrant`, `ThemeDark`, `ThemeNature`, `ThemeTech`, `ThemeCarbon`

### `AllThemes() []Theme`

Returns all built-in theme presets.

### `Theme` struct fields

```go
type Theme struct {
    Name       string
    Colors     ColorScheme
    Fonts      FontScheme
    Primary    string
    Secondary  string
    Accent     string
    Background string
    Text       string
    Light      string
    Dark       string
}
```

### `ColorScheme` struct fields

`Name`, `Dk1`, `Lt1`, `Dk2`, `Lt2`, `Accent1`–`Accent6`, `Hlink`, `FolHlink`

### `FontScheme` struct fields

`Name`, `MajorFont`, `MinorFont`

Pass a `Theme` to `PresentationBuilder.WithTheme(theme)`.

## See also

- [API Reference](../api-reference.md)
- [Python Presentation API](python-presentation-api.md)
- [JSON Bridge Operations](bridge-operations.md)
- [Go Text Reference](go-text.md)
- [Go Hyperlinks Reference](go-hyperlinks.md)
- [Go Export Reference](go-export.md)
- [Go Layout Helpers Reference](go-layout-helpers.md)
