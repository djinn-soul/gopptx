# I10 - Animation Effects

**Focus:** Add animations to slides for visual effects.

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func main() {
	pres := pptx.NewPresentationBuilder("I10 Animation")

	slide1 := pptx.NewSlide("Intro").AddBullet("Animated content")
	pres.AddSlide(slide1)

	slide2 := pptx.NewSlide("Main Point").
		AddBullet("First item").
		AddBullet("Second item").
		WithTransitionOptions(transitions.TransitionOptions{Type: transitions.TransitionFade, DurationMS: 500})
	pres.AddSlide(slide2)

	_ = pres.WriteToFile("i10-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.transitions import TRANSITION_FADE

with Presentation.new("I10 Animation") as p:
    p.add_slide("Intro")
    p.slides[0].add_textbox(0.8, 2.0, 8.0, 1.5, text="Animated content")

    p.add_slide("Main Point")
    slide2 = p.slides[1]
    slide2.add_textbox(0.8, 2.0, 8.0, 2.0, text="First item\nSecond item")
    slide2.set_transition(TRANSITION_FADE, duration_ms=500)
    p.save("docs/assets/pptx/usage/i10-python.pptx")
```

**Download PPTX:** [i10-python.pptx](../../../assets/pptx/usage/i10-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Animation Effects](../../../assets/images/usage/i10-python.png)
