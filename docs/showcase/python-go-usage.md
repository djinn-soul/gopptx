# Python + Go Usage (Same Workflow)

This page shows the same practical workflow in both Python and Go, producing equivalent output.

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

- `python/gopptx/presentation/runtime_lifecycle.py`
- `docs/guides/python-workflows.md`

## Go Usage

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	err := pptx.NewPresentationBuilder("Team Weekly Update").
		AddTitleSlide("Executive Summary").
		AddBulletSlide("Highlights", []string{
			"Revenue up 8%",
			"Retention improved by 4%",
			"Deployment lead time reduced",
		}).
		WriteToFile("team_weekly_go.pptx")
	if err != nil {
		panic(err)
	}
}
```

Go source references:

- `pkg/pptx/presentation_builder.go`
- `docs/guides/go-workflows.md`

## Output

Both examples produce a `.pptx` file with:

1. A title slide — **"Executive Summary"**
2. A bullet slide — **"Highlights"** with three bullet points

The Go and Python APIs are equivalent in structure: both use a builder-style flow to compose slides and write a single file at the end.
