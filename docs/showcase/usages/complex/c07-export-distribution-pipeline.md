# C07 - Export & Distribution Pipeline

**Focus:** Automate export and distribution workflows.

**Go code**

```go
package main

import (
	"os"
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
)

func main() {
	pres := pptx.NewPresentationBuilder("C07 Export Pipeline")
	pres.AddSlide(pptx.NewSlide("Export Demo").AddBullet("PDF export").AddBullet("HTML export"))

	pptxBytes, _ := pres.Build()

	os.WriteFile("c07.pptx", pptxBytes, 0644)

	_ = export.PDFFromFile("c07.pptx", "c07.pdf")
}
```

**Python code**

```python
import os
from gopptx import Presentation

with Presentation.new("C07 Export Pipeline") as p:
    p.add_slide("Export Demo")
    p.slides[0].add_textbox(0.8, 2.0, 8.0, 2.0, text="PDF export\nHTML export")

    p.save("c07.pptx")
    p.export_pdf("c07.pdf")
    p.export_html("c07")
```

**Download PPTX:** [c07-python.pptx](../../../assets/pptx/usage/c07-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Export Distribution Pipeline](../../../assets/images/usage/c07-python.png)
