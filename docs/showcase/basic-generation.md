# Basic Generation (Python)

Build a deck from scratch with one slide and one shape.

## Code

```python
from gopptx import Inches, Presentation

with Presentation.new("Basic Generation") as pres:
    pres.add_slide("Hello World")
    pres.add_shape(0, "rect", (Inches(1), Inches(1), Inches(2), Inches(1)), text="Basic Shape")
    pres.save("basic_generation.pptx")
```

## Run It

```bash
go run ./examples/01-basic-pptx-generation/basic_gen.go
```

## Artifacts

- Source: `examples/01-basic-pptx-generation/basic_gen.go`
- PPTX: [basic-generation.pptx](../assets/pptx/basic-generation.pptx)
- Screenshot:

![Basic Generation Screenshot](../assets/images/showcase/basic-gen.png)
