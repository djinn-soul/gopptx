# Quickstart

Two paths — Python and Go. Both produce a real `.pptx` you can open in PowerPoint.

!!! warning "Geometry is measured in EMU"

    English Metric Units: **914 400 per inch**, 12 700 per point. A shape placed at
    `(40, 120, 600, 220)` is six hundred EMU wide — about 0.0007 inch — and looks like it never
    got added. Always wrap coordinates:

    ```python
    from gopptx import Emu, Inches, Point
    ```

    In Go the same helpers are `pptx.Inches()`, `pptx.Points()`, `pptx.Centimeters()`,
    `pptx.Emu()`. Note that the convenience shape constructors (`pptx.NewRectangle` and friends)
    are the exception — they take **inches as plain floats**.

## Python

### 1. Create a deck

```python
from gopptx import Presentation

with Presentation.new("Quarterly Update") as pres:
    pres.add_bullet_slide("Highlights", ["Revenue +12%", "Retention +4%"])
    pres.add_paragraph_slide("Summary", "We beat plan in every region.")
    pres.save("quickstart.pptx")
```

`Presentation.new(title)` returns a deck that **already has a title slide at index 0**. The
`with` block owns the native handle and releases it on exit; `save()` is what writes bytes.

### 2. Add real content

```python
from gopptx import ChartType, Inches, Presentation, ShapeType

with Presentation.new("Content") as pres:
    slide = pres.add_slide("Numbers")

    pres.add_textbox(
        slide.index,
        Inches(0.5), Inches(0.4), Inches(4.0), Inches(0.6),
        text="Revenue by quarter",
    )

    pres.add_shape(
        slide.index,
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(0.5), Inches(1.2), Inches(2.0), Inches(0.8)),
        text="On plan",
    )

    pres.add_table_from_rows(
        slide.index,
        [["Region", "Revenue"], ["EMEA", "4.1M"], ["APAC", "2.8M"]],
        bounds=(Inches(0.5), Inches(2.4), Inches(4.0), Inches(1.8)),
    )

    pres.add_chart(
        slide.index,
        ChartType.COLUMN,
        ["Q1", "Q2", "Q3", "Q4"],
        [12.0, 15.5, 18.0, 21.0],
        bounds=(Inches(5.0), Inches(1.2), Inches(4.0), Inches(3.0)),
        title="Quarterly trend",
    )

    pres.save("content.pptx")
```

### 3. Edit an existing deck

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.set_slide_title(0, "Updated Title")
    pres.find_and_replace("Draft", "Final")
    pres.add_slide("Closing")
    pres.save("edited.pptx")
```

`Presentation(path)`, `Presentation.open_deck(path)` and `pres.open(path)` all open an existing
file. Use whichever reads best; they are the same operation.

### 4. Navigate with the object layer

```python
with Presentation("edited.pptx") as pres:
    for slide in pres.slides:
        print(slide.index, slide.title, len(slide.list_shapes()))

    first = pres.slides[0]
    first.notes = "Open with the revenue number."
    first.set_transition("fade")
    pres.save("annotated.pptx")
```

### 5. Export

```python
from gopptx import HTMLOptions, PDFOptions, Presentation

with Presentation("edited.pptx") as pres:
    pres.export_pdf("deck.pdf", PDFOptions(driver="auto"))
    pres.export_html("deck.html", HTMLOptions(embed_images=True))
```

See the [export guide](guides/export.md) for the driver trade-offs.

## Go

### 1. Create a deck in one call

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

	data, err := pptx.CreateWithSlides("Quickstart Deck", slides)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("quickstart.pptx", data, 0o600); err != nil {
		panic(err)
	}
}
```

### 2. Build a deck fluently

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	agenda := pptx.NewSlideBuilder("Agenda")
	agenda.AddBullet("Results")
	agenda.AddBullet("Outlook")
	agenda.WithNotes("Two minutes, no more.")

	err := pptx.NewPresentationBuilder("Quarterly Update").
		WithTheme(pptx.ThemeCorporate).
		WithSlideSize(pptx.SlideSize16x9()).
		AddTitleSlide("FY26 Q3").
		AddSlide(agenda.Build()).
		WriteToFile("quarterly.pptx")
	if err != nil {
		panic(err)
	}
}
```

!!! tip "Prefer `SlideBuilder` over `SlideContent`"

    `SlideContent`'s chainable methods use a **value receiver**, so dropping the return value
    discards the change and the compiler stays silent:

    ```go
    s := pptx.NewSlide("t")
    s.AddBullet("lost")   // no effect — return value discarded
    ```

    `SlideBuilder` mutates in place, so the same mistake cannot happen:

    ```go
    b := pptx.NewSlideBuilder("t")
    b.AddBullet("kept")   // applied
    slide := b.Build()
    ```

### 3. Add shapes and a chart

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	b := pptx.NewSlideBuilder("Numbers")

	// Convenience constructors take inches as plain floats.
	b.AddShape(
		shapes.NewRoundedRectangle(0.5, 1.2, 2.0, 0.8).
			WithText("On plan").
			WithFill(shapes.NewShapeFill("DCE6F2")).
			WithLine(shapes.NewShapeLine("1F4E78", styling.Points(1))),
	)

	b.WithBarChart(pptx.NewBarChart(
		[]string{"Q1", "Q2", "Q3", "Q4"},
		[]float64{12, 15.5, 18, 21},
	))

	if err := pptx.NewPresentationBuilder("Numbers").
		AddSlide(b.Build()).
		WriteToFile("numbers.pptx"); err != nil {
		panic(err)
	}
}
```

### 4. Edit an existing deck

```go
p, err := pptx.Open("input.pptx")
if err != nil {
	panic(err)
}
defer p.Close()

p.SetTitle("Updated Title")
p.SetAuthor("Reporting Bot")

if err := p.SaveAs("edited.pptx"); err != nil {
	panic(err)
}
```

`pptx.Open` gives you document properties and chart data. For shape-level editing — the same
surface the JSON bridge drives — use `pptx.OpenEditor(path)`, which returns a
`*PresentationEditor`.

## Next steps

- [Core concepts](concepts.md) — handles, envelopes, units, batching
- [Python library guide](guides/python-library.md) / [Go library guide](guides/go-library.md)
- [Showcase](showcase/usages/index.md) — 30 annotated recipes with screenshots
- [`examples/`](https://github.com/djinn-soul/gopptx/tree/main/examples) — 96 runnable programs
