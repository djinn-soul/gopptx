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
- `Edit() (*Presentation, error)` — hand off straight into the editing API
- `WriteToFile(path string) error`

!!! tip "Prefer `SlideBuilder` when composing across statements"

    `SlideContent`'s chainable methods take a value receiver, so a dropped return value silently
    discards the change. `pptx.NewSlideBuilder(title)` mutates in place. See the
    [Go library guide](../guides/go-library.md#slidebuilder-vs-slidecontent).

## Styling: Length units

Source file: `pkg/pptx/styling/units.go`

All geometry is in EMU — 914 400 per inch. See [Units](../concepts.md#units-and-geometry).

- `Inches(value float64) Length`
- `Centimeters(value float64) Length`
- `Points(value float64) Length`
- `Emu(value int64) Length`
- `FontSize(pt float64) int`
- `(l Length).Inches() float64`
- `(l Length).Cm() float64`
- `(l Length).Pt() float64`
- `(l Length).Emu() int64`

!!! warning

    The convenience shape constructors — `shapes.NewRectangle`, `NewEllipse` and the rest —
    take **inches as plain `float64`**, not `Length`. `pptx.NewShape(prst, x, y, cx, cy)` takes
    `Length`.

### Relative dimensions

A length may also be stated against the slide:

- `Absolute(length Length) Dimension`
- `PercentOf(percent float64) Dimension`
- `Ratio(fraction float64) Dimension`

### Font size presets

`FontSizeTitle` (44), `FontSizeSubtitle` (32), `FontSizeHeading` (28), `FontSizeBody` (18), `FontSizeSmall` (14), `FontSizeCaption` (12), `FontSizeCode` (14), `FontSizeLarge` (36), `FontSizeXLarge` (48)

## Styling: Colors

Source file: `pkg/pptx/styling/colors.go`

Basic: `ColorRed`, `ColorGreen`, `ColorBlue`, `ColorWhite`, `ColorBlack`, `ColorGray`, `ColorLightGray`, `ColorDarkGray`, `ColorYellow`, `ColorLightBlue`, `ColorOrange`, `ColorPurple`, `ColorCyan`, `ColorMagenta`, `ColorNavy`, `ColorTeal`, `ColorOlive`

Corporate: `ColorCorporateBlue`, `ColorCorporateGreen`, `ColorCorporateRed`, `ColorCorporateOrange`

Material Design: `ColorMaterialRed`, `ColorMaterialPink`, `ColorMaterialPurple`, `ColorMaterialIndigo`, `ColorMaterialBlue`, `ColorMaterialCyan`, `ColorMaterialTeal`, `ColorMaterialGreen`, `ColorMaterialLime`, `ColorMaterialAmber`, `ColorMaterialOrange`, `ColorMaterialBrown`, `ColorMaterialGray`

IBM Carbon: `ColorCarbonBlue60`, `ColorCarbonBlue40`, `ColorCarbonGray100`, `ColorCarbonGray80`, `ColorCarbonGray20`, `ColorCarbonGreen50`, `ColorCarbonRed60`, `ColorCarbonPurple60`

### `ColorValue` — colour arithmetic

The constants above are hex strings. `ColorValue` is a real colour type with the arithmetic
callers used to do by hand. The zero value is opaque black, so build one with a constructor
rather than a struct literal:

```go
c := pptx.MustColorFromHex("1F4E78")
c = pptx.RGB(31, 78, 120)
c = pptx.RGBA(31, 78, 120, 128)
```

It carries `Darker`, `Lighter`, blending and `ContrastRatio` — useful for deriving an accessible
text colour from a brand colour rather than guessing.

## Styling: Line dash constants

Source file: `pkg/pptx/styling/line_style.go`

`LineDashSolid` (`solid`), `LineDashDash` (`dash`), `LineDashDot` (`dot`),
`LineDashDashDot` (`dashDot`), `LineDashDashDotDot` (`lgDashDotDot`),
`LineDashLongDash` (`lgDash`), `LineDashLongDashDot` (`lgDashDot`),
`LineDashLongDashDotDot` (`lgDashDotDot`), `LineDashSystemDash` (`sysDash`),
`LineDashSystemDot` (`sysDot`), `LineDashSystemDashDot` (`sysDashDot`)

!!! note

    `ST_PresetLineDashVal` has no short dash-dot-dot token, so `LineDashDashDotDot` and
    `LineDashLongDashDotDot` both map to `lgDashDotDot`. That is deliberate, not a duplicate.

## Styling: Built-in themes

Source file: `pkg/pptx/styling/theme.go`

### Theme variables

`ThemeCorporate`, `ThemeModern`, `ThemeVibrant`, `ThemeDark`, `ThemeNature`, `ThemeTech`,
`ThemeCarbon`, `ThemeOffice`

### `ResolveTheme(name string) (Theme, bool)`

Resolves a theme name to a preset. Accepts the gopptx names above and the Office preset names
(`office`, `office2013`, `facet`, `integral`, `ion`, `retrospect`, `slice`, `wisp`), ignoring
case and separators. `ThemeNames()` lists every accepted name. Both `apply_theme` and
`set_global_theme_preset` resolve through it, so the two share one vocabulary.

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

## Keeping this page honest

This page summarises; `go doc` is generated from the source and cannot drift:

```bash
go doc github.com/djinn-soul/gopptx/pkg/pptx
go doc github.com/djinn-soul/gopptx/pkg/pptx.PresentationBuilder
go doc github.com/djinn-soul/gopptx/pkg/pptx/elements.SlideBuilder
```

## See also

- [Go library guide](../guides/go-library.md)
- [API Reference](../api-reference.md)
- [Feature matrix](feature-matrix.md)
- [Python Presentation API](python-presentation-api.md)
- [JSON Bridge Operations](bridge-operations.md)
- [Go Text Reference](go-text.md)
- [Go Hyperlinks Reference](go-hyperlinks.md)
- [Go Export Reference](go-export.md)
- [Go Layout Helpers Reference](go-layout-helpers.md)
