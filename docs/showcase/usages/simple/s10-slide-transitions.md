# S10 - Slide Transitions

**Focus:** Apply visual transitions between slides (Fade, Morph, etc.).

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/gopptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/transitions"
)

func main() {
	pres := &gopptx.Presentation{Title: "S10 Slide Transitions"}
	pres.AddSlide().Title = "Slide 1: Start"

	slide2 := pres.AddSlide()
	slide2.Title = "Slide 2: Fade Entry"
	slide2.WithTransitionOptions(transitions.TransitionOptions{Type: transitions.TransitionFade, DurationMS: 1000})

	_ = pres.Save("s10-go.pptx")
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
