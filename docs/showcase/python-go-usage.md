# Python + Go Usage (Same Workflow)

This page shows the same practical workflow in both Python and Go, with matching output screenshots.

## Workflow

1. Create a presentation.
2. Add a structured content slide.
3. Export `.pptx`.
4. Review rendered slide output.

## Python Usage

```python
from gopptx import Presentation

with Presentation.new("Team Weekly Update") as pres:
    pres.add_slide("Executive Summary")
    pres.add_bullet_slide(
        "Highlights",
        [
            "Revenue up 8%",
            "Retention improved by 4%",
            "Deployment lead time reduced"
        ],
    )
    pres.save("team_weekly_python.pptx")
```

Python source references:

- `python/gopptx/presentation/runtime.py`
- `docs/guides/python-workflows.md`

## Go Usage

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	p := pptx.NewPresentation()
	s := p.AddSlide()
	s.AddTextBox("Executive Summary")
	s.AddTextBox("Revenue up 8% | Retention +4% | Faster delivery")
	_ = p.Save("team_weekly_go.pptx")
}
```

Go source references:

- `pkg/pptx`
- `docs/guides/go-workflows.md`

## Output Screenshots

### Basic Deck Output

![Basic generation](../assets/images/showcase/basic-gen.png)

### Rich Slide Output

![Rich slide](../assets/images/showcase/rich-slide.png)

### Brand + Theme Output

![Brand reskin](../assets/images/showcase/reskin-result.png)

### Chart Output

![Radar chart](../assets/images/showcase/chart-radar.png)

## Download PPTX Samples

- [basic-generation.pptx](../assets/pptx/basic-generation.pptx)
- [rich-slide.pptx](../assets/pptx/rich-slide.pptx)
- [brand-reskin.pptx](../assets/pptx/brand-reskin.pptx)
- [chart-radar.pptx](../assets/pptx/chart-radar.pptx)
