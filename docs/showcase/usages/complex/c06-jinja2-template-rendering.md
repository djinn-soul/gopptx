# C06 - Jinja2 Template Rendering

**Focus:** Use Jinja2 templates for dynamic content.

**Go code**

Go uses [gonja](https://github.com/noirbizarre/gonja) (a Go-native Jinja2 engine) via the `render_template` bridge operation.
The workflow is: build the PPTX with Jinja2 template text in slide shapes, then call `render_template` on the editor to expand it in place.

```go
package main

import (
	"encoding/json"
	"os"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
)

func main() {
	// Step 1: build a deck with Jinja2 template expressions in slide text.
	data, err := pptx.NewPresentationBuilder("C06 Jinja2").
		AddSlide(
			pptx.NewSlide("{{ title }}").
				AddBullet("{% for section in sections %}- {{ section.title }}: {{ section.content }}{% endfor %}"),
		).
		Build()
	if err != nil {
		panic(err)
	}

	// Step 2: open the raw bytes with the editor.
	ed, err := editor.OpenPresentationEditorFromBytes(data)
	if err != nil {
		panic(err)
	}
	defer ed.Close()

	// Step 3: render all Jinja2 expressions in place via the bridge operation.
	ctx := map[string]any{
		"title": "Jinja2 Report",
		"sections": []map[string]any{
			{"title": "Overview", "content": "Summary here"},
			{"title": "Details", "content": "More info"},
		},
	}
	payload, _ := json.Marshal(map[string]any{
		"api_version": 1,
		"op":          "render_template",
		"payload":     map[string]any{"context": ctx},
	})
	editor.ExecuteCommand(ed, string(payload))

	// Step 4: save.
	out, err := ed.SaveToBytes()
	if err != nil {
		panic(err)
	}
	_ = os.WriteFile("c06-go.pptx", out, 0o600)
}
```

**Python code**

```python
from jinja2 import Template
from gopptx import Presentation

template = Template("""
Title: {{ title }}
{% for section in sections %}
- {{ section.title }}: {{ section.content }}
{% endfor %}
""")

data = {
    "title": "Jinja2 Report",
    "sections": [
        {"title": "Overview", "content": "Summary here"},
        {"title": "Details", "content": "More info"},
    ],
}

with Presentation.new("C06 Jinja2") as p:
    content = template.render(**data)
    p.add_slide("Rendered")
    p.slides[0].add_textbox(0.8, 2.0, 8.0, 3.0, text=content)
    p.save("docs/assets/pptx/usage/c06-python.pptx")
```

**Download PPTX:** [c06-python.pptx](../../../assets/pptx/usage/c06-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Jinja2 Template Rendering](../../../assets/images/usage/c06-python.png)
