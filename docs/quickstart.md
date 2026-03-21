# Quickstart

This quickstart gives you one Python path and one Go path.

## Python: Create a New Deck

```python
from gopptx import Presentation

with Presentation.new("Quarterly Update") as pres:
    pres.add_slide("Overview")
    pres.add_bullet_slide("Highlights", ["Growth +12%", "Retention +4%"])
    pres.save("quickstart_python.pptx")
```

## Python: Edit an Existing Deck

```python
from gopptx import Presentation

with Presentation("input.pptx") as pres:
    pres.set_slide_title(0, "Updated Title")
    pres.add_slide("New Closing Slide")
    pres.save("quickstart_python_edited.pptx")
```

## Go: Minimal Deck

```bash
go get github.com/djinn-soul/gopptx
```

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	p := pptx.NewPresentation()
	slide := p.AddSlide()
	slide.AddTextBox("Hello from gopptx")
	_ = p.Save("quickstart_go.pptx")
}
```

## Next Step

1. Browse runnable examples with screenshots in [Examples](showcase/index.md).
2. Read [Python Library](guides/python-library.md) or [Go Library](guides/go-library.md).
3. Use [Reference](api-reference.md) for full method/operation coverage.