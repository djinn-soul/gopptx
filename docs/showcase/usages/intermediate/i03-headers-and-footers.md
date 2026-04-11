# I03 - Headers and Footers

**Focus:** Add headers and footers to slides.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	slide := pptx.NewSlide("Headers and Footers").
		AddBullet("Headers provide context").
		AddBullet("Footers add metadata")

	_ = pptx.NewPresentationBuilder("I03 Headers and Footers").
		AddSlide(slide).
		WriteToFile("i03-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation

with Presentation.new("I03 Headers and Footers") as p:
    slide = p.add_slide("Headers and Footers")
    slide.header = "Company Presentation"
    slide.footer = "Confidential - 2023"
    slide.add_textbox(0.8, 2.0, 8.0, 2.0, text="Headers provide context\nFooters add metadata")
    p.save("docs/assets/pptx/usage/i03-python.pptx")
```

**Download PPTX:** [i03-python.pptx](../../../assets/pptx/usage/i03-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Headers and Footers](../../../assets/images/usage/i03-python.png)
