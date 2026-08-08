# Go library guide

```bash
go get github.com/djinn-soul/gopptx
```

```go
import "github.com/djinn-soul/gopptx/pkg/pptx"
```

`pkg/pptx` is a façade: it re-exports the types from `pkg/pptx/shapes`, `charts`, `tables`,
`text`, `styling`, `elements`, `templates` and friends, so most programs need only this one
import. Reach into the sub-packages when you want the narrower names.

## The three entry points

| You want to… | Use | Returns |
|---|---|---|
| Emit a deck in one call | `pptx.CreateWithSlides`, `pptx.Create`, `pptx.WriteFile` | `[]byte` / file |
| Compose a deck fluently | `pptx.NewPresentationBuilder(title)` | `*PresentationBuilder` |
| Open and edit an existing deck | `pptx.Open` / `pptx.OpenEditor` | `*Presentation` / `*PresentationEditor` |

## Generating a deck

```go
package main

import (
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Hello from gopptx").AddBullet("Created with gopptx"),
	}

	data, err := pptx.CreateWithSlides("My Deck", slides)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("output.pptx", data, 0o600); err != nil {
		panic(err)
	}
}
```

`pptx.WriteFile(path, title, slides)` does the same thing straight to disk.

## SlideBuilder vs SlideContent

!!! danger "The most common Go mistake in this library"

    `SlideContent`'s chainable methods take a **value receiver**. Dropping the return value
    silently discards the change and the compiler says nothing:

    ```go
    s := pptx.NewSlide("Agenda")
    s.AddBullet("lost")        // no effect
    ```

    Chaining inline is fine, because nothing is dropped:

    ```go
    s := pptx.NewSlide("Agenda").AddBullet("kept").AddBullet("also kept")
    ```

    `SlideBuilder` mutates in place, so the mistake is impossible. Prefer it whenever a slide is
    built across more than one statement:

    ```go
    b := pptx.NewSlideBuilder("Agenda")
    b.AddBullet("kept")
    b.WithNotes("Two minutes.")
    slide := b.Build()
    ```

    `pptx.BuildFrom(existing)` wraps a `SlideContent` you already have.

## Composing with PresentationBuilder

```go
agenda := pptx.NewSlideBuilder("Agenda")
agenda.AddBullet("Results")
agenda.AddSubBullet(1, "By region")
agenda.AddNumbered("Then outlook")
agenda.WithNotes("Two minutes, no more.")

// pptx.Metadata embeds common.Metadata (aliased as pptx.MetadataFields),
// so the document fields are easiest to set on a zero value.
var meta pptx.Metadata
meta.Title = "FY26 Q3"
meta.Creator = "Reporting Bot"

err := pptx.NewPresentationBuilder("Quarterly Update").
	WithTheme(pptx.ThemeCorporate).
	WithSlideSize(pptx.SlideSize16x9()).
	WithMetadata(meta).
	WithSlideNumbers(true).
	WithFooter("Confidential").
	AddTitleSlide("FY26 Q3").
	AddSlide(agenda.Build()).
	WriteToFile("quarterly.pptx")
```

Other terminators: `Build() ([]byte, error)` for bytes, and `Edit() (*Presentation, error)` to
carry straight on into the editing API.

Built-in themes are real values here — `pptx.ThemeCorporate`, `ThemeModern`, `ThemeDark`,
`ThemeVibrant`, `ThemeNature`, `ThemeTech`, `ThemeCarbon`, `ThemeOffice`, and
`pptx.AllThemes()`.

## Units

```go
pptx.Inches(1.5)       // styling.Length (EMU)
pptx.Centimeters(2.0)
pptx.Points(12)
pptx.Emu(914400)
```

Relative dimensions, where a length is stated against the slide:

```go
pptx.Absolute(pptx.Inches(2))
pptx.PercentOf(50)     // half the slide
pptx.Ratio(0.25)
```

!!! note "The convenience shape constructors take inches, not `Length`"

    `shapes.NewRectangle(1.0, 1.8, 2.2, 1.3)` is inches as plain `float64`.
    `pptx.NewShape("rect", pptx.Inches(1.0), …)` is EMU `Length`. Mixing them up produces
    shapes that are either invisible or larger than the slide.

## Shapes

```go
import (
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

b := pptx.NewSlideBuilder("Numbers")

b.AddShape(
	shapes.NewRoundedRectangle(0.5, 1.2, 2.0, 0.8).
		WithText("On plan").
		WithFill(shapes.NewShapeFill("DCE6F2")).
		WithLine(shapes.NewShapeLine("1F4E78", styling.Points(1))),
)

b.AddConnector(
	shapes.NewElbowConnector(
		styling.Inches(2.6), styling.Inches(2.7),
		styling.Inches(4.6), styling.Inches(2.7),
	).WithLine(shapes.NewShapeLine("444444", styling.Points(1.5))),
)

b.AutoRerouteConnectors()
```

There are roughly 200 preset constructors — `NewEllipse`, `NewChevron`, `NewCube`, the whole
`NewFlowChart*` family, `NewStar4` through `NewStar32`, the callouts and the action buttons.
For anything else, `pptx.NewShape(prstName, x, y, cx, cy)` takes the raw OOXML preset name and
the `pptx.ShapeType*` constants enumerate the valid ones.

Richer fills and effects:

```go
shapes.NewShapeGradientFill("linear", []shapes.ShapeGradientStop{
	shapes.NewShapeGradientStop(0, "1F4E78"),
	shapes.NewShapeGradientStop(100, "DCE6F2"),
})
pptx.NewOuterShadow("404040")
```

Freeforms:

```go
ff, err := pptx.NewFreeformInches([][2]float64{{1, 1}, {2, 2}, {1, 2}})
```

## Text

```go
b.AddBullet("Plain bullet")
b.AddSubBullet(1, "Indented one level")
b.AddNumbered("Numbered item")
b.AddLettered("Lettered item")
b.AddFormattedBullet("Supports **bold** and *italic* inline")

b.AddBulletRuns([]pptx.Run{
	pptx.NewRun("Revenue "),
	pptx.NewRun("+12%"),
})
```

Paragraph-level control comes from `pptx.NewParagraphStyle()` and the `AddBulletWithStyle` /
`AddBulletRunsWithStyle` variants.

## Tables

```go
t := pptx.NewTable([]pptx.Length{pptx.Inches(2), pptx.Inches(2)})
// populate rows, then attach it to a slide placeholder:
b.WithPlaceholderTable(1, t)
```

## Charts

```go
b.WithBarChart(pptx.NewBarChart(
	[]string{"Q1", "Q2", "Q3", "Q4"},
	[]float64{12, 15.5, 18, 21},
))
```

Every chart family has a constructor: `NewLineChart`, `NewLineMarkersChart`, `NewPieChart`,
`NewDoughnutChart`, `NewAreaStackedChart`, `NewRadarFilledChart`, `NewScatterChart`,
`NewBubbleChart`, `NewStockOHLCChart`, `NewComboChart`, and so on — 24 in total. The generic
form is `pptx.NewChart(enums.XLChartType…, categories, values)`.

For multi-series data use the builders:

```go
data := pptx.NewCategoryChartData([]string{"Q1", "Q2"})
// add series to data, then hand it to the chart
```

A generated chart embeds its Excel workbook, so the numbers remain editable in PowerPoint.

## Layout helpers

Rather than hand-computing offsets:

```go
boxes, err := pptx.Grid(2, 3, pptx.Inches(0.25))                 // 2×3 grid over the slide
boxes, err = pptx.GridInBox(2, 3, pptx.Inches(0.25), bounds)     // …within a box

points, err := pptx.Stack("horizontal", start, pptx.Inches(0.2), sizes...)
offsets, err := pptx.DistributeUniform("horizontal", bounds, 4, pptx.Inches(1.5))
offsets, err = pptx.DistributeNonUniform("vertical", bounds, sizes)

x, y := pptx.Center(pptx.Inches(2), pptx.Inches(1))
x, y = pptx.CenterInBox(pptx.Inches(2), pptx.Inches(1), bounds)
```

## Editing an existing deck

Two different objects, for two different jobs.

### `pptx.Open` — properties and chart data

```go
p, err := pptx.Open("input.pptx")
if err != nil {
	return err
}
defer p.Close()

p.SetTitle("Updated Title")
p.SetAuthor("Reporting Bot")
p.SetKeywords("qbr,finance")

refs, _ := p.ListSlideCharts(2)
_ = p.ReplaceChartData(2, 0, []string{"Q1", "Q2"}, []float64{10, 14})

issues := p.Validate()
return p.SaveAs("edited.pptx")
```

`OpenFromBytes` and `OpenFromReader` are the in-memory equivalents; `SaveToBytes` and
`SaveToWriter` complete the pair. `Save()` writes back over the file it was opened from.

### `pptx.OpenEditor` — the full shape tree

```go
import "github.com/djinn-soul/gopptx/pkg/pptx/editor"

e, err := pptx.OpenEditor("input.pptx")
if err != nil {
	return err
}

resp := editor.ExecuteCommand(e, `{
	"api_version": 1,
	"op": "add_slide",
	"payload": {"title": "Appendix"}
}`)   // resp is the JSON response envelope, as a string

return pptx.Save(e, "edited.pptx")
```

`PresentationEditor` is the same object the JSON bridge drives, so anything Python can do is
reachable here — all 179 operations. See
[Bridge operations](../reference/bridge-operations.md).

## Importing content

```go
slides, err := pptx.SlidesFromMarkdown("# Agenda\n- One\n- Two")
slides, err = pptx.SlidesFromMarkdownFile("deck.md")
```

Web content goes through `pkg/pptx/urlfetch`:

```go
data, err := urlfetch.URLToPPTX("https://example.com")
```

## Validation and repair

```go
issues, err := pptx.Validate(data)
fixed, result, err := pptx.Repair(data)
```

Validation covers OPC package rules, required parts, relationship targets and content types —
the defects that make PowerPoint offer to repair a file. It does not catch a deck that is
structurally valid but visually wrong; open the result before trusting a layout.

## Examples

96 runnable programs live under `examples/`:

```bash
go run ./examples/01-basic-pptx-generation
go run ./examples/12-shapes
go run ./examples/58-gopptx-rich-slide
```

Useful starting points:

| Example | Shows |
|---|---|
| `01-basic-pptx-generation` | The minimum viable deck |
| `12-shapes` | Preset geometry, fills, lines, effects |
| `09-charts` | Every chart family |
| `19-read-modify-existing` | Read-modify-save |
| `24-smartart` | SmartArt authoring |
| `32-mermaid` | Mermaid diagrams |
| `35-layout-helpers` | Grid, stack, distribute |
| `58-gopptx-rich-slide` | Everything at once |
| `81-native-pdf-showcase` | Native PDF export |

## Reference

- [Go API reference](../reference/go-api.md)
- `go doc github.com/djinn-soul/gopptx/pkg/pptx` — always current
