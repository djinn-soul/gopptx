# I04 - Grouped Content Slide

**Focus:** Combine text, image, and shapes in one slide.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	slide := pptx.NewSlide("").WithBlankLayout().
		AddShape(pptx.NewTextBox("Grouped Content Slide", 0.8, 0.35, 5.8, 0.45)).
		AddShape(
			pptx.NewRoundedRectangle(0.75, 1.2, 3.35, 4.75).
				WithFill(pptx.NewShapeFill("EAF2FC")).
				WithLine(pptx.NewShapeLine("A9C4E2", pptx.Points(1.2))).
				WithText("Text panel\n\nText block for context\nImage for visual anchor\nShapes for emphasis").
				WithAutoFit(pptx.TextAutoFitNormal),
		).
		AddShape(
			pptx.NewRoundedRectangle(4.35, 1.2, 4.1, 3.05).
				WithFill(pptx.NewShapeFill("FFFFFF")).
				WithLine(pptx.NewShapeLine("C9D3E0", pptx.Points(1.1))),
		).
		AddImage(
			pptx.NewImage(
				"examples/assets/test_image.png",
				pptx.Inches(4.55),
				pptx.Inches(1.4),
				pptx.Inches(3.7),
				pptx.Inches(2.65),
			),
		).
		AddShape(
			pptx.NewRoundedRectangle(4.35, 4.45, 1.15, 0.45).
				WithText("Text").
				WithFill(pptx.NewShapeFill("D9E1F2")).
				WithLine(pptx.NewShapeLine("D9E1F2", pptx.Points(1.0))).
				WithAutoFit(pptx.TextAutoFitNormal),
		).
		AddShape(
			pptx.NewRoundedRectangle(5.65, 4.45, 1.15, 0.45).
				WithText("Image").
				WithFill(pptx.NewShapeFill("E2F0D9")).
				WithLine(pptx.NewShapeLine("E2F0D9", pptx.Points(1.0))).
				WithAutoFit(pptx.TextAutoFitNormal),
		).
		AddShape(
			pptx.NewRoundedRectangle(6.95, 4.45, 1.15, 0.45).
				WithText("Shapes").
				WithFill(pptx.NewShapeFill("FCE4D6")).
				WithLine(pptx.NewShapeLine("FCE4D6", pptx.Points(1.0))).
				WithAutoFit(pptx.TextAutoFitNormal),
		)
	pres := pptx.NewPresentationBuilder("I04 Grouped Content Slide")
	pres.AddSlide(slide)
	_ = pres.WriteToFile("i04-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation, ShapeType
from gopptx.schemas import Inches

with Presentation.new("I04 Grouped Content Slide") as p:
    p.update_slide(0, layout="blank")
    p.add_textbox(
        0,
        Inches(0.8),
        Inches(0.35),
        Inches(5.8),
        Inches(0.45),
        text="Grouped Content Slide",
    )
    p.add_shape(
        0,
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(0.75), Inches(1.2), Inches(3.35), Inches(4.75)),
        text="Text panel\n\nText block for context\nImage for visual anchor\nShapes for emphasis",
        properties={
            "fill": {"solid": "EAF2FC"},
            "line": {"color": "A9C4E2", "width_emu": 12700},
        },
    )
    p.add_shape(
        0,
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(4.35), Inches(1.2), Inches(4.1), Inches(3.05)),
        properties={
            "fill": {"solid": "FFFFFF"},
            "line": {"color": "C9D3E0", "width_emu": 12700},
        },
    )
    p.add_image(
        0,
        "examples/assets/test_image.png",
        (Inches(4.55), Inches(1.4), Inches(3.7), Inches(2.65)),
    )
    p.add_shape(
        0,
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(4.35), Inches(4.45), Inches(1.15), Inches(0.45)),
        text="Text",
        properties={
            "fill": {"solid": "D9E1F2"},
            "line": {"color": "D9E1F2", "width_emu": 12700},
        },
    )
    p.add_shape(
        0,
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(5.65), Inches(4.45), Inches(1.15), Inches(0.45)),
        text="Image",
        properties={
            "fill": {"solid": "E2F0D9"},
            "line": {"color": "E2F0D9", "width_emu": 12700},
        },
    )
    p.add_shape(
        0,
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(6.95), Inches(4.45), Inches(1.15), Inches(0.45)),
        text="Shapes",
        properties={
            "fill": {"solid": "FCE4D6"},
            "line": {"color": "FCE4D6", "width_emu": 12700},
        },
    )
    p.save("docs/assets/pptx/usage/i04-python.pptx")
```

**Download PPTX:** [i04-python.pptx](../../../assets/pptx/usage/i04-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Grouped Content Slide](../../../assets/images/usage/i04-python.png)
