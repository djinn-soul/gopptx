# C05 - Template + Data Injection System

**Focus:** Use templates with dynamic data injection.

**Go code**

```go
package main

import (
	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/templates"
)

type ReportData struct {
	Title    string
	Period   string
	Revenue  int
	Expenses int
	Profit   int
}

func main() {
	data := ReportData{
		Title:    "Q4 Financial Report",
		Period:   "Q4 2023",
		Revenue:  1250000,
		Expenses: 890000,
		Profit:   360000,
	}

	template := templates.LoadTemplate("financial_report")
	result := template.Execute(data)

	_ = pptx.NewPresentationBuilder("C05 Template").
		AddSlide(result).
		WriteToFile("c05-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation

with Presentation.new("C05 Template") as p:
    data = {
        "title": "Q4 Financial Report",
        "period": "Q4 2023",
        "revenue": 1250000,
        "expenses": 890000,
        "profit": 360000,
    }
    p.apply_template("financial", data)
    p.save("docs/assets/pptx/usage/c05-python.pptx")
```

**Download PPTX:** [c05-python.pptx](../../../assets/pptx/usage/c05-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Template + Data Injection System](../../../assets/images/usage/c05-python.png)
