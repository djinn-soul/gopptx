# S08 - Create a Blank Presentation

**Focus:** Start a new PPTX file and save it without visible slide content.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	_ = pptx.NewPresentationBuilder("Blank Presentation").
		AddSlide(
			pptx.NewSlide("").WithBlankLayout(),
		).
		WriteToFile("s08-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation

with Presentation.new("Blank Presentation") as p:
    p.add_slide("", layout="blank")
    p.remove_slide(0)
    p.save("docs/assets/pptx/usage/s08-python.pptx")
```

**Download PPTX:** [s08-python.pptx](../../../assets/pptx/usage/s08-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Create a Blank Presentation](../../../assets/images/usage/s08-python.png)
