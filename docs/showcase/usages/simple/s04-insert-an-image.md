# S04 - Insert an Image

**Focus:** Add a PNG/JPG image to a slide.

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/gopptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	pres := &gopptx.Presentation{Title: "S04 Insert an Image"}
	slide := pres.AddSlide()
	slide.Title = "Insert an Image"
	slide.AddImage(
		shapes.NewImage(
			"examples/assets/55/repository-open-graph-template.png",
			styling.Inches(0.8),
			styling.Inches(1.4),
			styling.Inches(8.0),
			styling.Inches(4.6),
		).WithAltText("Inserted PNG sample"),
	)
	_ = pres.Save("s04-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.schemas import Inches

with Presentation.new("S04 Insert an Image") as p:
    p.add_image(
        0,
        "examples/assets/55/repository-open-graph-template.png",
        (Inches(0.8), Inches(1.4), Inches(8.0), Inches(4.6)),
    )
    p.save("docs/assets/pptx/usage/s04-python.pptx")
```

**Download PPTX:** [s04-python.pptx](../../../assets/pptx/usage/s04-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Insert an Image](../../../assets/images/usage/s04-python.png)
