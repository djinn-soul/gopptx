# S05 - Set Slide Background Color

**Focus:** Apply a solid background color.

**Go code**

```go
package main

import (
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	slide := pptx.NewSlide("Slide Background Color").
		WithBackgroundColor("D9E1F2").
		AddBullet("Apply a solid background color.").
		AddBullet("Keep title and content readable.").
		AddBullet("Save as PPTX.")

	data, err := pptx.CreateWithSlides("S05 Set Slide Background Color", []pptx.SlideContent{slide})
	if err != nil {
		panic(err)
	}
	_ = os.WriteFile("s05-go.pptx", data, 0o600)
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.constants import ConnectorType, ShapeType
from gopptx.schemas import Inches

with Presentation.new("S05 Set Slide Background Color") as p:
    slide = p.slides[0]
    slide.title = "Slide Background Color"
    slide.set_background("solid", color="D9E1F2")
    slide.add_textbox(
        Inches(0.8),
        Inches(2.0),
        Inches(8.0),
        Inches(1.6),
        text="Solid background color applied: D9E1F2",
    )
    p.save("docs/assets/pptx/usage/s05-python.pptx")
```

**Download PPTX:** [s05-python.pptx](../../../assets/pptx/usage/s05-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Set Slide Background Color](../../../assets/images/usage/s05-python.png)
