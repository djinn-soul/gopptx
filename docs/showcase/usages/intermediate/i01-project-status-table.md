# I01 - Project Status Table

**Focus:** Present workstreams, owners, and status in a simple table.

**Go code**

```go
package main

import "github.com/djinn-soul/gopptx/pkg/pptx"

func main() {
	pres := pptx.NewPresentationBuilder("I01 Project Status Table")
	pres.AddSlide(
		pptx.NewSlide("Project Status Table").WithTable(
			pptx.NewTable([]pptx.Length{pptx.Inches(2.6), pptx.Inches(2.6), pptx.Inches(2.6)}).
				WithStyledData([][]pptx.TableCell{
					{
						pptx.NewTableCell("Workstream").WithBold(true).WithBackgroundColor("DCE6F2"),
						pptx.NewTableCell("Owner").WithBold(true).WithBackgroundColor("DCE6F2"),
						pptx.NewTableCell("Status").WithBold(true).WithBackgroundColor("DCE6F2"),
					},
					{
						pptx.NewTableCell("Docs revamp"),
						pptx.NewTableCell("Design Ops"),
						pptx.NewTableCell("Ready"),
					},
					{
						pptx.NewTableCell("Usage catalog"),
						pptx.NewTableCell("Docs Team"),
						pptx.NewTableCell("Queued"),
					},
				}),
		),
	)
	_ = pres.WriteToFile("i01-go.pptx")
}
```

**Python code**

```python
from gopptx import Presentation
from gopptx.schemas import Inches

with Presentation.new("I01 Project Status Table") as p:
    p.add_slide("Project Status Table")
    p.add_table(
        0,
        [
            ["Workstream", "Owner", "Status"],
            ["Docs revamp", "Design Ops", "Ready"],
            ["Usage catalog", "Docs Team", "Queued"],
        ],
        column_widths=[Inches(2.6), Inches(2.6), Inches(2.6)],
        header=True,
    )
    p.save("docs/assets/pptx/usage/i01-python.pptx")
```

**Download PPTX:** [i01-python.pptx](../../../assets/pptx/usage/i01-python.pptx)

Screenshot generated from the Python code above using `export_pptx_png.ps1`.

![Project Status Table](../../../assets/images/usage/i01-python.png)
