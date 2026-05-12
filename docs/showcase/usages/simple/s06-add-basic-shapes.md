# S06 - Add Basic Shapes

**Focus:** Insert rectangle, circle, and line shapes.

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	slide := pptx.NewSlide("Basic Shapes").
		AddShape(
			shapes.NewRectangle(1.0, 1.8, 2.2, 1.3).
				WithText("Rectangle").
				WithFill(shapes.NewShapeFill("DCE6F2")).
				WithLine(shapes.NewShapeLine("1F4E78", styling.Points(1))),
		).
		AddShape(
			shapes.NewEllipse(4.6, 1.8, 1.8, 1.8).
				WithText("Circle").
				WithFill(shapes.NewShapeFill("FCE4D6")).
				WithLine(shapes.NewShapeLine("9C3F00", styling.Points(1))),
		).
		AddConnector(
			shapes.NewStraightConnector(
				styling.Inches(2.6),
				styling.Inches(2.7),
				styling.Inches(4.6),
				styling.Inches(2.7),
			).WithLine(shapes.NewShapeLine("444444", styling.Points(1.5))),
		)
	_ = pptx.NewPresentationBuilder("S06 Add Basic Shapes").AddSlide(slide).WriteToFile("s06-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.constants import ConnectorType, ShapeType
from gopptx.schemas import Inches

with Presentation.new("S06 Add Basic Shapes") as p:
    slide = p.slides[0]
    slide.title = "Basic Shapes"
    slide.add_shape(
        ShapeType.RECTANGLE,
        (Inches(1.0), Inches(1.8), Inches(2.2), Inches(1.3)),
        text="Rectangle",
        properties={
            "fill": {"solid": "DCE6F2"},
            "line": {"color": "1F4E78", "width_emu": 12700},
        },
    )
    slide.add_shape(
        ShapeType.ELLIPSE,
        (Inches(4.6), Inches(1.8), Inches(1.8), Inches(1.8)),
        text="Circle",
        properties={
            "fill": {"solid": "FCE4D6"},
            "line": {"color": "9C3F00", "width_emu": 12700},
        },
    )
    slide.add_connector(
        ConnectorType.STRAIGHT,
        Inches(2.6),
        Inches(2.7),
        Inches(4.6),
        Inches(2.7),
        properties={"line": {"color": "444444", "width_emu": 19050}},
    )
    p.save("docs/assets/pptx/usage/s06-python.pptx")
```

**Download PPTX:** [s06-python.pptx](../../../assets/pptx/usage/s06-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Add Basic Shapes](../../../assets/images/usage/s06-python.png)
