# C08 - Narrative Storyboard Generator

**Focus:** Generate storyboard presentations from narrative content.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

type StoryScene struct {
	Title   string
	Content string
	Visual  string
}

func main() {
	scenes := []StoryScene{
		{"Opening", "The journey begins", "scene1.png"},
		{"Rising Action", "Challenge emerges", "scene2.png"},
		{"Climax", "Turning point", "scene3.png"},
		{"Resolution", "Outcome achieved", "scene4.png"},
	}

	pres := pptx.NewPresentationBuilder("C08 Storyboard")

	for _, scene := range scenes {
		slide := pptx.NewSlide(scene.Title).
			AddBullet(scene.Content)
		pres.AddSlide(slide)
	}

	_ = pres.WriteToFile("c08-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation

scenes = [
    {"title": "Opening", "content": "The journey begins", "visual": "scene1.png"},
    {"title": "Rising Action", "content": "Challenge emerges", "visual": "scene2.png"},
    {"title": "Climax", "content": "Turning point", "visual": "scene3.png"},
    {"title": "Resolution", "content": "Outcome achieved", "visual": "scene4.png"},
]

with Presentation.new("C08 Storyboard") as p:
    for scene in scenes:
        p.add_slide(scene["title"])
        p.slides[-1].add_textbox(0.8, 2.0, 8.0, 2.0, text=scene["content"])
    p.save("docs/assets/pptx/usage/c08-python.pptx")
```

**Download PPTX:** [c08-python.pptx](../../../assets/pptx/usage/c08-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Narrative Storyboard Generator](../../../assets/images/usage/c08-python.png)
