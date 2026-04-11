# S03 - Paragraph Text

**Focus:** Add a readable paragraph text block.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	builder := pptx.NewPresentationBuilder("S03 Paragraph Text")
	s := pptx.NewSlide("Paragraph Text").
		AddShape(
			pptx.NewTextBox(
				"Paragraph text helps explain context clearly. "+
					"Use one concise block for intent, then keep follow-up details compact.",
				0.8, 2.0, 8.0, 2.2,
			),
		).
		AddShape(
			pptx.NewTextBox(
				"Add another paragraph on the same slide using a second text box block.",
				0.8, 4.5, 8.0, 1.0,
			),
		)
	_ = builder.AddSlide(s).WriteToFile("s03-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation

with Presentation.new("S03 Paragraph Text") as p:
    slide = p.add_paragraph_slide(
        "Paragraph Text",
        (
            "Paragraph text helps explain context clearly. "
            "Use one concise block for intent, then keep follow-up details compact."
        ),
    )
    slide.add_paragraph(
        "Add another paragraph on the same slide without manually setting textbox bounds."
    )
    p.save("docs/assets/pptx/usage/s03-python.pptx")
```

**Download PPTX:** [s03-python.pptx](../../../assets/pptx/usage/s03-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Paragraph Text](../../../assets/images/usage/s03-paragraph-python.png)
