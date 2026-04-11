# C06 - Jinja2 Template Rendering

**Focus:** Use Jinja2 templates for dynamic content.

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/templates"
)

func main() {
	template := `
Title: {{ title }}
{% for section in sections %}
- {{ section.title }}: {{ section.content }}
{% endfor %}
`
	data := map[string]interface{}{
		"title": "Jinja2 Report",
		"sections": []map[string]string{
			{"title": "Overview", "content": "Summary here"},
			{"title": "Details", "content": "More info"},
		},
	}

	result := templates.RenderJinja2(template, data)
	_ = pptx.NewPresentationBuilder("C06 Jinja2").
		AddSlide(pptx.NewSlide("Rendered").WithContent(result)).
		WriteToFile("c06-go.pptx")
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
