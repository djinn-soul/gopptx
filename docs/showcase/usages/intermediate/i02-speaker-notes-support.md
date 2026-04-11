# I02 - Speaker Notes Support

**Focus:** Add speaker notes to slides for presentation guidance.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	slide := pptx.NewSlide("Speaker Notes Support").
		AddBullet("Present with confidence").
		AddBullet("Use notes for key points").
		WithNotes("Remember to emphasize the benefits of automation.")

	_ = pptx.NewPresentationBuilder("I02 Speaker Notes Support").
		AddSlide(slide).
		WriteToFile("i02-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation

with Presentation.new("I02 Speaker Notes Support") as p:
    slide = p.add_slide("Speaker Notes Support")
    slide.add_textbox(0.8, 2.0, 8.0, 2.0, text="Present with confidence\nUse notes for key points")
    slide.notes = "Remember to emphasize the benefits of automation."
    p.save("docs/assets/pptx/usage/i02-python.pptx")
```

**Download PPTX:** [i02-python.pptx](../../../assets/pptx/usage/i02-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Speaker Notes Support](../../../assets/images/usage/i02-python.png)
