# S07 - Styled Text Formatting

**Focus:** Use font size, bold, italic, underline, and color.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	slide := pptx.NewSlide("Styled Text Formatting").
		WithTitleColor("1F4E79").
		WithTitleBold(true).
		WithTitleItalic(true).
		WithContentSize(22).
		WithContentBold(true).
		WithContentItalic(true).
		WithContentUnderline(true).
		WithContentColor("C00000").
		AddBullet("Formatted bullet content")

	_ = pptx.NewPresentationBuilder("S07 Styled Text Formatting").
		AddSlide(slide).
		WriteToFile("s07-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.schemas import Inches

with Presentation.new("S07 Styled Text Formatting") as p:
    slide = p.add_slide("Styled Text Formatting")
    slide.add_textbox(
        Inches(0.8),
        Inches(1.8),
        Inches(8.0),
        Inches(2.0),
        runs=[
            {"text": "Font size 34pt  ", "size_pt": 34},
            {"text": "Bold  ", "bold": True, "size_pt": 28},
            {"text": "Italic  ", "italic": True, "size_pt": 28},
            {"text": "Underline  ", "underline": "sng", "size_pt": 28},
            {"text": "Color", "color": "C00000", "size_pt": 28},
        ],
    )
    p.save("docs/assets/pptx/usage/s07-python.pptx")
```

**Download PPTX:** [s07-python.pptx](../../../assets/pptx/usage/s07-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Styled Text Formatting](../../../assets/images/usage/s07-python.png)
