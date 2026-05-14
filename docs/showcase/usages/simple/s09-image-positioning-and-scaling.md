# S09 - Image Positioning and Scaling

**Focus:** Resize and place images precisely.

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	slide := pptx.NewSlide("").
		AddImage(
			shapes.NewImage(
				"examples/assets/test_image.png",
				styling.Inches(0.5),
				styling.Inches(1.0),
				styling.Inches(3.0),
				styling.Inches(3.0),
			),
		).
		AddImage(
			shapes.NewImage(
				"examples/assets/test_image.png",
				styling.Inches(4.0),
				styling.Inches(1.0),
				styling.Inches(2.2),
				styling.Inches(2.2),
			),
		).
		AddImage(
			shapes.NewImage(
				"examples/assets/test_image.png",
				styling.Inches(6.7),
				styling.Inches(2.0),
				styling.Inches(1.4),
				styling.Inches(1.4),
			),
		)
	_ = pptx.NewPresentationBuilder("S09 Image Positioning and Scaling").AddSlide(slide).WriteToFile("s09-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.schemas import Inches

with Presentation.new("S09 Image Positioning and Scaling") as p:
    p.update_slide(0, layout="blank")
    slide = p.slides[0]
    slide.add_image(
        "examples/assets/test_image.png",
        (Inches(0.5), Inches(1.0), Inches(3.0), Inches(3.0)),
    )
    slide.add_image(
        "examples/assets/test_image.png",
        (Inches(4.0), Inches(1.0), Inches(2.2), Inches(2.2)),
    )
    slide.add_image(
        "examples/assets/test_image.png",
        (Inches(6.7), Inches(2.0), Inches(1.4), Inches(1.4)),
    )
    p.save("docs/assets/pptx/usage/s09-python.pptx")
```

**Download PPTX:** [s09-python.pptx](../../../assets/pptx/usage/s09-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Image Positioning and Scaling](../../../assets/images/usage/s09-python.png)
