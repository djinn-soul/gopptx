# C09 - SmartArt Layouts

**Focus:** Use SmartArt for visual content layouts.

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func main() {
	pres := pptx.NewPresentationBuilder("C09 SmartArt")

	slide := pptx.NewSlide("SmartArt Layout").
		AddSmartArt(
			smartart.NewSmartArt(smartart.BasicBlockList).
				AddItems([]string{
					"Process step 1",
					"Process step 2",
					"Process step 3",
					"Process step 4",
				}),
		)

	pres.AddSlide(slide)
	_ = pres.WriteToFile("c09-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.constants import SMARTART_PROCESS

with Presentation.new("C09 SmartArt") as p:
    p.add_slide("SmartArt Layout")
    p.slides[0].add_smartart(
        SMARTART_PROCESS,
        ["Process step 1", "Process step 2", "Process step 3", "Process step 4"],
    )
    p.save("docs/assets/pptx/usage/c09-python.pptx")
```

**Download PPTX:** [c09-python.pptx](../../../assets/pptx/usage/c09-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![SmartArt Layouts](../../../assets/images/usage/c09-python.png)
