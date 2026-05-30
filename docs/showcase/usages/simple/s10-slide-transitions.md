# S10 - Slide Transitions

**Focus:** Apply visual transitions between slides (Fade, Morph, etc.).

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func main() {
	slide1 := pptx.NewSlide("Slide 1: Start")
	slide2 := pptx.NewSlide("Slide 2: Fade Entry").
		WithTransitionOptions(transitions.TransitionOptions{Type: transitions.TransitionFade, DurationMS: 1000})

	_ = pptx.NewPresentationBuilder("S10 Slide Transitions").
		AddSlide(slide1).
		AddSlide(slide2).
		WriteToFile("s10-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.transitions import TRANSITION_FADE

with Presentation.new("S10 Slide Transitions") as p:
    p.add_slide("Slide 1: Start")
    slide2 = p.add_slide("Slide 2: Fade Entry")
    slide2.set_transition(TRANSITION_FADE, duration_ms=1000)
    p.save("docs/assets/pptx/usage/s10-python.pptx")
```

**Download PPTX:** [s10-python.pptx](../../../assets/pptx/usage/s10-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Slide Transitions](../../../assets/images/usage/s10-python.png)
