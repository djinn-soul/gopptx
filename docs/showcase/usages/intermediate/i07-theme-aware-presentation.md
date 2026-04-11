# I07 - Theme-aware Presentation

**Focus:** Apply themes consistently across slides.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	pres := pptx.NewPresentationBuilder("I07 Theme Aware").
		WithTheme(pptx.ThemeCorporate)

	pres.AddSlide(pptx.NewSlide("Title Slide").
		AddBullet("Company").
		AddBullet("Q4 Results"))

	pres.AddSlide(pptx.NewSlide("Agenda").
		AddBullet("Revenue growth").
		AddBullet("Cost reduction").
		AddBullet("Future plans"))

	_ = pres.WriteToFile("i07-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation

with Presentation.new("I07 Theme Aware") as p:
    p.apply_theme("corporate")
    p.add_slide("Title Slide")
    p.slides[0].title = "Company Q4 Results"
    p.add_slide("Agenda")
    p.slides[1].add_textbox(0.8, 2.0, 8.0, 2.0, text="Revenue growth\nCost reduction\nFuture plans")
    p.save("docs/assets/pptx/usage/i07-python.pptx")
```

**Download PPTX:** [i07-python.pptx](../../../assets/pptx/usage/i07-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Theme-aware Presentation](../../../assets/images/usage/i07-python.png)
