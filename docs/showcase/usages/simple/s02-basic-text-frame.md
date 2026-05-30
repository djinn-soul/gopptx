# S02 - Basic Text Frame

**Focus:** Add controlled text regions.

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

func main() {
	slide := pptx.NewSlide("Basic Text Frame").
		AddShape(shapes.NewRectangle(0.8, 2.0, 3.0, 1.0).WithText("Top anchor sample")).
		AddShape(shapes.NewRectangle(4.4, 2.0, 3.0, 1.0).WithText("Bottom anchor sample")).
		AddShape(shapes.NewRectangle(0.8, 3.3, 6.6, 1.0).WithText("No-wrap / shrink-fit text region"))
	_ = pptx.NewPresentationBuilder("S02 Text Frame").AddSlide(slide).WriteToFile("s02-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.schemas import Inches, Point

with Presentation.new("S02 Text Frame") as p:
    p.add_slide("Basic Text Frame")
    p.add_shape(
        0,
        "rect",
        (Inches(0.8), Inches(1.6), Inches(3.0), Inches(1.0)),
        text="Top anchor sample",
        properties={"fill": {"solid": "FFF2CC"}, "line": {"color": "B45F06", "width_emu": Point(2)}},
    )
    p.add_shape(
        0,
        "rect",
        (Inches(4.2), Inches(1.6), Inches(3.0), Inches(1.0)),
        text="Bottom anchor sample",
        properties={"fill": {"solid": "D9EAD3"}, "line": {"color": "38761D", "width_emu": Point(2)}},
    )
    p.add_shape(
        0,
        "rect",
        (Inches(0.8), Inches(3.0), Inches(6.4), Inches(1.0)),
        text="No-wrap / shrink-fit text region",
        properties={"fill": {"solid": "D0E0E3"}, "line": {"color": "134F5C", "width_emu": Point(2)}},
    )
    p.save("docs/assets/pptx/usage/s02-python.pptx")
```

**Download PPTX:** [s02-python.pptx](../../../assets/pptx/usage/s02-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Basic Text Frame](../../../assets/images/usage/s02-python.png)
